/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/aokabi/gogit/pkg"
	"github.com/aokabi/gogit/pkg/config"
	"github.com/spf13/cobra"
)

const (
// GIT_HOST = "http://hostos:8084/git"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push repo branch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("need 2 args")
			return
		}

		repo := args[0]
		branch := args[1]
		conf := config.Read()
		remoteUrl := conf.GetRemoteUrl(repo)
		refspec := fmt.Sprintf("refs/heads/%s", branch)

		// とりあえずHTTPプロトコルで実装してみる
		// SSHより扱いに慣れてるから

		resp, err := getGitReceivePack(remoteUrl)
		if err != nil {
			fmt.Println(err)
			return
		}

		// ローカルの状態を取得
		newHash := strings.Trim(pkg.ReadRef(refspec), "\n")

		// capも決め打ち
		// 公式gitの通信内容を見て真似している
		r := newReceivePackRequest(resp.refs[refspec], newHash, refspec, "report-status-v2 side-band-64k object-format=sha1", createPackfile())
		fmt.Println(r)
		if err := postGitReceivePack(&r, remoteUrl); err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type getReceivePackResponse struct {
	refs map[string]string //key: ref, value: sha
	caps []string
}

func New() *getReceivePackResponse {
	return &getReceivePackResponse{
		refs: map[string]string{},
		caps: []string{},
	}
}

func (s *getReceivePackResponse) String() string {
	return fmt.Sprintf("refs: %v, caps: %v", s.refs, s.caps)
}

func getGitReceivePack(remoteUrl string) (*getReceivePackResponse, error) {
	apiUrl, err := url.JoinPath(remoteUrl, "/info/refs")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	query := req.URL.Query()
	query.Add("service", "git-receive-pack")
	req.URL.RawQuery = query.Encode()
	// 本当はGIT_ASKPASSのコマンドを使ったりしてちゃんとしたい
	req.SetBasicAuth(os.Getenv("GITHUB_ID"), os.Getenv("GITHUB_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// parse response
	// cf. https://git-scm.com/docs/gitprotocol-http#_smart_clients

	respBody := New()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	first := scanner.Text()
	validFirst := regexp.MustCompile(`^[0-9a-f]{4}#`)
	if !validFirst.MatchString(first) {
		return nil, fmt.Errorf("invalid first line: %s", first)
	}

	processing := false
	for scanner.Scan() {
		line := scanner.Text()

		if line[0:4] == "0000" {
			if processing {
				break
			} else {
				processing = true
				line = line[4:]
			}
		}

		_ = line[0:4] // length
		line = line[4:]
		if strings.Contains(line, pkg.NullByte) {
			tmp := strings.Split(line, pkg.NullByte)
			line = tmp[0]

			caps := strings.Split(tmp[1], " ")
			respBody.caps = append(respBody.caps, caps...)
		}

		refs := strings.Split(line, " ")
		respBody.refs[refs[1]] = refs[0]
	}

	return respBody, nil
}

// cf. https://git-scm.com/docs/http-protocol#_smart_service_git_receive_pack
func postGitReceivePack(r *receivePackRequest, remoteUrl string) error {
	apiUrl, err := url.JoinPath(remoteUrl, "/git-receive-pack")
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPost, apiUrl, r.newRequestBody())

	// 本当はGIT_ASKPASSのコマンドを使ったりしてちゃんとしたい
	req.SetBasicAuth(os.Getenv("GITHUB_ID"), os.Getenv("GITHUB_TOKEN"))

	req.Header.Set("Content-Type", "application/x-git-receive-pack-request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	fmt.Println(string(respBody))

	return nil

}

type receivePackRequest struct {
	length  int
	oldHash string
	newHash string

	refs string
	cap  string

	packfileData []byte
}

func (r receivePackRequest) String() string {
	return fmt.Sprintf("%d %s %s %s %s", r.length, r.oldHash, r.newHash, r.refs, r.cap)
}

func newReceivePackRequest(oldHash string, newHash string, refs string, cap string, packfileData []byte) receivePackRequest {
	return receivePackRequest{
		// 4 + はlength自身, 最後の+1は改行分
		length:  4 + len(oldHash) + 1 + len(newHash) + 1 + len(refs) + 1 + len(cap) + 1,
		oldHash: oldHash,
		newHash: newHash,
		refs:    refs,
		cap:     cap,

		packfileData: packfileData,
	}
}

func (r receivePackRequest) newRequestBody() io.Reader {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("%04x%s %s %s\x00 %s", r.length, r.oldHash, r.newHash, r.refs, r.cap))
	buf.WriteString("0000")
	fmt.Println(buf.String())
	buf.Write(r.packfileData)

	return buf
}

func createPackfile() []byte {
	// exec external command
	// 今後これも自前で実装したい
	// あと現状すべてが含まれたでかいpackfileができてしまう
	// gitの挙動を見てるとサーバー側との差分だけpackfileにしていそうなので、そうしたい
	packObjects := exec.Command("git", "pack-objects", "--stdout", "--revs", "--delta-base-offset", "--thin")
	stdin, _ := packObjects.StdinPipe()
	head := pkg.ReadRef(pkg.ReadHEAD())

	func() {
		defer stdin.Close()
		if _, err := io.WriteString(stdin, head); err != nil {
			panic(err)
		}
	}()

	out, err := packObjects.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return out

}
