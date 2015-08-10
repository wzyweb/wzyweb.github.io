package main
import (
        "os"
        "fmt"
        "strings"
        "strconv"
        "regexp"
)

var i int64=0
var lujing string=""
var newname string
var lastname string
var mode string=""
var modeFile string=""
var text string
var img string
var title string
var wri string=""
var post string
var time string=""
var tag string=""
var tagurl string=""
var tagurlpre string=""
var nav string=""
var body string=""

func main() {
    readtext()
    creattagurl()
    newfilename()
    fmt.Printf("%s\n",newname)
    readmode()
    bodytext()
    texttohtml()
    writefile()
}

func newfilename() {
       lujing = "./post/"
       dir, _ := os.Open(lujing)
       files, _ := dir.Readdir(0)
       var b string
       var c []string
       for _, a := range files {
           if !a.IsDir() {
               b = a.Name()
               c = strings.Split(b,"-")
               d, err := strconv.ParseInt(c[0],10,64)
               if err != nil {
                   panic(err)
               }
               if i<d {
                   i = d
                   tagurlpre = strings.Split(c[1],".")[0]
               }
           }
       }
       newname = strconv.FormatInt(i+1,10) + "-" + tagurl + ".htm"
       lastname = strconv.FormatInt(i,10) + "-" + tagurlpre + ".htm"
}

func readmode() {
    fin,err := os.Open(modeFile)
    defer fin.Close()
    if err != nil {
             fmt.Println(modeFile,err)
             return
     }
     buf := make([]byte, 1024)
     for{
            n, _ := fin.Read(buf)
            if 0 == n { break }
            list := string(buf[:n])
            mode = mode + list
     }

}

func readtext() {
    textFile := "./post.md"
    fin,err := os.Open(textFile)
    defer fin.Close()
    if err != nil {
             fmt.Println(textFile,err)
             return
     }
     buf := make([]byte, 1024)
     alltext := ""
     for{
            n, _ := fin.Read(buf)
            if 0 == n { break }
            alltext = alltext + string(buf[:n])
     }
            sp := strings.Split(alltext,"\r\n\r\n")
            sp2 := strings.Split(sp[0],"\r\n")
            if strings.Contains(sp2[0],"[img-H]") {
                img = strings.Replace(sp2[0],"[img-H]","",-1)
                modeFile = "./mo/article.mo"
            }
            if strings.Contains(sp2[0],"[img-S]") {
                img = strings.Replace(sp2[0],"[img-S]","",-1)
                modeFile = "./mo/article2.mo"
            }
            time = sp2[1]
            tag = sp2[2]
            temp := sp[0]+"\r\n\r\n"+sp[1]+"\r\n\r\n"
            text = strings.Replace(alltext,temp,"",1)
            title = strings.Replace(sp[1],"## ","",1)
            fmt.Printf("tag: %s\n",tag)
            fmt.Printf("img's URL is: %s\n",img)

}

func creattagurl() {
    if tag == "照片随笔" {
        tagurl = "photostory"
    }
    if tag == "七种武器" {
        tagurl = "equipment"
    }
    

}


func texttohtml() {
    meta := "</div><div class=\"meta\">Written on " + time + " | Tag: <a href=\"../archive/" + tagurl + "1.htm\">" + tag + "</a>"
    message := "</div><div class=\"mess\"><a href=\"http://www.douban.com/group/heimaphoto/\" target=\"_blank\">留言</a>"
    post = "<h2>" + title + "</h2>" + body + meta + message
    if i == 0 {
        nav = ""
    } else {
        nav = "<a href=\"./" + lastname + "\">上一篇</a>"
    }
}

func writefile() {
    wFile := lujing + newname
    fout,err := os.Create(wFile)
    defer fout.Close()
    if err != nil {
            fmt.Println(wFile,err)
            return
    }
    wri = strings.Replace(mode,"[title]",title,1)
    wri = strings.Replace(wri,"[img]",img,1)
    wri = strings.Replace(wri,"[post]",post,1)
    wri = strings.Replace(wri,"[nav]",nav,1)
    fout.WriteString(wri)
    if i != 0 {
        preFile := lujing + lastname
        prefout, er := os.Open(preFile)
        defer prefout.Close()
        if er != nil {
            fmt.Println(preFile,er)
            return
        }
        buf := make([]byte, 1024)
        alltext := ""
        for{
            n, _ := prefout.Read(buf)
            if 0 == n { break }
            alltext = alltext + string(buf[:n])
        }
        if !strings.Contains(alltext,"下一篇</a>") {
        if i == 1 {
            oldtext := "索引页</a> | "
            replacetext := oldtext + "<a href=\"./" + newname + "\">下一篇</a>"
            alltext = strings.Replace(alltext,oldtext,replacetext,1)
            prewr, eerr := os.Create(preFile)
            defer prewr.Close()
            if eerr != nil {
                fmt.Println(preFile,eerr)
                return
            }
            prewr.WriteString(alltext)
        } else {
            oldtext := "上一篇</a> | "
            replacetext := oldtext + " | <a href=\"./" + newname + "\">下一篇</a>"
            alltext = strings.Replace(alltext,oldtext,replacetext,1)
            prewr, eerr := os.Create(preFile)
            defer prewr.Close()
            if eerr != nil {
                fmt.Println(preFile,eerr)
                return
            }
            prewr.WriteString(alltext)
        }
        }
    }
}

func bodytext() {
    body = "<p>"+strings.Replace(text,"\r\n\r\n","</p>\r\n<p>",-1)+"</p>\r\n"
    body = strings.Replace(body," \r\n","<br />\r\n",-1)
    reg := regexp.MustCompile(`\+ +|- +`)
    list := reg.FindAllString(body,-1)
    st2 := ""
    if len(list) > 0 {
        bodylist := strings.Split(body,"</p>")
        reg2 := regexp.MustCompile(`<p>\+ +|<p>- +`)
        st := ""
        st3 := ""
        for k, v := range bodylist {
           switch k {
                case len(bodylist)-1:
                default:
                if len(reg2.FindAllString(v,-1)) > 0 {
                    if strings.Contains(v,"\r\n") {
                        reg3 := regexp.MustCompile(`\r\n\+ +|\r\n- +`)
                        st3 = strings.Replace(v,reg3.FindAllString(v,-1)[0],"</li><li>",-1)
                        st = strings.Replace(st3,reg2.FindAllString(v,-1)[0],"<p><li>",1)
                        st2 = st2 + st + "</li></p>"
                    } else {
                    st = strings.Replace(v,reg2.FindAllString(v,-1)[0],"<p><li>",1)
                    st2 = st2 + st + "</li></p>"
                    }
                } else {
                    st2 = st2 + v + "</p>"
                }
            }
        }
        body = st2
    }
}
