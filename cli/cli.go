package cli

import (
	"bufio"
	"io"
	"strings"
	"text/tabwriter"
)

// RenderReplaceHeader rendering but with header replaced
func RenderReplaceHeader(input io.Reader, output io.Writer, flags uint, headers []string) error {
	w := tabwriter.NewWriter(output, 0, 0, 3, ' ', flags)
	defer w.Flush()

	buf := bufio.NewReader(input)
	_, _, err := buf.ReadLine()
	if err != nil {
		return err
	}

	totalColumnHeader := len(headers)
	bodyByte := []byte{}
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		j := strings.Split(string(line), ",")
		totalColumnBody := len(j)
		for i := 0; i < totalColumnHeader; i++ {
			if totalColumnBody >= totalColumnHeader {
				// kalo misalkan column body lebih banyak dari pada column header maka column body di cut
				bodyByte = append(bodyByte, []byte(j[i])...)
				if i < totalColumnHeader-1 {
					bodyByte = append(bodyByte, []byte("\t")...)
				}
			} else {
				// kalo misalkan column header lebih banyak dari pada column body maka column body di tambahin
				if i >= totalColumnBody {
					bodyByte = append(bodyByte, []byte("")...)
				} else {
					bodyByte = append(bodyByte, []byte(j[i])...)
				}
				if i < totalColumnHeader-1 {
					bodyByte = append(bodyByte, []byte("\t")...)
				}
			}
		}
		bodyByte = append(bodyByte, []byte("\n")...)
	}

	header := strings.Join(headers, "\t") + "\n"
	// hapus new line terakhir
	body := strings.TrimSuffix(string(bodyByte), "\n")
	all := header + body

	text := strings.ReplaceAll(all, ",", "\t")

	_, err = w.Write([]byte(text))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}
