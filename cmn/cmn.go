package cmn

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const (
	DbCreate = iota
	DbRead
	DbUpdate
	DbDelete
)

// Reads a file to string
func ReadFileToString(filename string) (string, error) {
	var b []byte

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Gets a string between two tags
func StringBetweenTags(src string, openTag string, closeTag string) string {
	//https://github.com/google/re2/wiki/Syntax
	var re = regexp.MustCompile(`(?ms)` + openTag + `\n(.+)?\n` + closeTag)

	res := re.FindStringSubmatch(src)

	if len(res) > 1 {
		return res[1]
	} else {
		return ""
	}
}

// Gets a string between two tags
func StringBetweenTagsMap(src string) map[string]string {

	//clean carriage return
	src = strings.ReplaceAll(src, "\r", "")

	var str = make(map[string]string)

	//önce tagleri alalım
	var re = regexp.MustCompile(`(?ms)-- [a-z]+ >>`)
	res := re.FindAllString(src, -1)
	tags := make([]string, 0)
	tmpNdx := 0
	for _, tag := range res {
		tag = strings.ReplaceAll(tag, "-", "")
		tag = strings.ReplaceAll(tag, ">", "")
		tag = strings.TrimSpace(tag)

		if tmpNdx == 0 {
			tags = append(tags, tag)
			tmpNdx += 1
		} else {
			if tag != tags[tmpNdx-1] {
				tags = append(tags, tag)
				tmpNdx += 1
			}
		}
	}

	for _, tag := range tags {
		xre := regexp.MustCompile(fmt.Sprintf(`(?m)(?s)-- %s >>\n(.+)?\n-- %s >>`, tag, tag))
		xres := xre.FindStringSubmatch(src)
		if len(xres) > 1 {
			str[tag] = xres[0]
		}
	}

	return str
}

// Gets a string template between two tags
func TemplateBetweenTags(src string, openTag string, closeTag string) (*template.Template, error) {
	tstr := StringBetweenTags(src, openTag, closeTag)

	tmpl, err := template.New("test").Parse(tstr)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Int64InSlice(val int64, list []int64) bool {
	for _, sliceVal := range list {
		if sliceVal == val {
			return true
		}
	}
	return false
}

func StringSliceToInt64Slice(src []string) ([]int64, error) {
	var dst []int64

	for _, srcVal := range src {
		dstVal, err := strconv.ParseInt(srcVal, 10, 0)
		if err != nil {
			return nil, err
		}
		dst = append(dst, dstVal)
	}

	return dst, nil
}

func IntersectionInt64(a []int64, b []int64) (inter []int64) {
	// interacting on the smallest list first can potentailly be faster...but not by much, worse case is the same
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			// get future index values
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					// if the future values aren't the same then that's the end of the intersection
					if low[f1] != high[f2] {
						done = true
					}
				}
				// we don't want to interate on the entire list everytime, so remove the parts we already looped on will make it faster each pass
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		// nothing in the future so we are done
		if done {
			break
		}
	}
	return
}

var ModulPrefixes = []string{"cr", "df", "dm", "mm", "og", "fi", "sd", "pc", "pm", "pp", "pd", "qm", "rm", "ec", "ac"}

func FieldNameToModelName(fieldName string) string {
	varName := strcase.ToCamel(strings.ToLower(fieldName))
	if strings.HasSuffix(varName, "Id") {
		varName = varName[:len(varName)-2] + "ID"
	}

	prefix := strings.ToLower(fieldName[:2])
	if StringInSlice(prefix, ModulPrefixes) && (strings.HasSuffix(fieldName, "_id") || strings.HasSuffix(fieldName, "_code") || strings.HasSuffix(fieldName, "_title")) {
		varName = varName[:2] + strings.ToUpper(varName[2:3]) + varName[3:]
	}

	switch fieldName {
	case "crn":
		varName = "Crn"
	case "crcodetyp_id":
		varName = "CrCodeTypID"
	case "crcodetyp_code":
		varName = "CrCodeTypCode"
	}

	return varName
}

func ModelNameToFieldName(modelName string) string {

	switch modelName {
	case "CrCodeTypID":
		return "crcodetyp_id"
	case "CrCodeTypCode":
		return "crcodetyp_code"
	case "IpGroupCode":
		return "ipgroup_code"

	default:
		prefix := strings.ToLower(modelName[:2])
		if StringInSlice(prefix, ModulPrefixes) {
			modelName = strings.ToLower(modelName[:3]) + modelName[3:]
		}
		tmpStr := strcase.ToSnake(modelName)

		/*
			tablerequest'te cover ettiğim için kapattım
			if strings.HasSuffix(tmpStr, "_code") {
				res := strings.Split(tmpStr, "_")
				tmpStr = strings.Join(res, ".")
			}
			if strings.HasSuffix(tmpStr, "_title") {
				res := strings.Split(tmpStr, "_")
				tmpStr = strings.Join(res, ".")
			}
		*/

		return tmpStr
	}
}

func PrependSlice(dest []interface{}, arg interface{}) []interface{} {
	dest = append(dest, nil)
	copy(dest[1:], dest)
	dest[0] = arg
	return dest
}
