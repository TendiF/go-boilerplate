package utils

import (
  "regexp"
  "mime/multipart"
  "math/rand"
  "time"
  "strings"
  "fmt"
  "strconv"
  "github.com/gin-gonic/gin"
  "os"
  "../model"
  "../database"
 )

func CheckEmail(email string) bool {
  re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
  return re.MatchString(email)
}

func Pagination(item_per_page int, current_page int, total_item int ) map[string]int {
  res := make(map[string]int)
  if item_per_page == 0 {
    item_per_page = 5
  }
  if current_page == 0 {
    current_page = 1
  }
  res["current_page"] = current_page
  res["total_item"] = total_item
  // res["item_per_page"] = item_per_page

  res["offset"] = (current_page * item_per_page) - item_per_page;
  res["limit"] = item_per_page;
  return res;
}

func Upload(file *multipart.FileHeader, c *gin.Context)int{
  created_by, err := c.Get("user_id")
  u_type := "user"
  if !err {
    created_by = c.MustGet("admin_id")
    u_type = "admin"
  }
  // single file
  t := time.Now()
  if _, err := os.Stat("./public"); os.IsNotExist(err) {
    dErr := os.Mkdir("./public", 0755)
    if dErr != nil {
      gin.DefaultWriter.Write([]byte("ERROR : "+dErr.Error()))
    }
  }
  dest := "./public/uploads/"
  if _, err := os.Stat(dest); os.IsNotExist(err) {
    dErr := os.Mkdir(dest, 0755)
    if dErr != nil {
      gin.DefaultWriter.Write([]byte("ERROR : "+dErr.Error()))
    }
  }
  current_time := t.Format("20060102150405")
  rand_number := rand.Intn(9999)
  slice_data := strings.Split(file.Filename, ".")
  format_file := slice_data[len(slice_data) - 1]
  file_name := strings.Replace(file.Filename, "."+format_file, "", -1)
  strByte := RandStringBytes(5)
  dest = dest +file_name + current_time + strByte + strconv.Itoa(rand_number) + "." + format_file
  // Upload the file to specific dst.
  res := c.SaveUploadedFile(file, dest)
  fmt.Println(res)
  return model.Add("uploads", map[string]string{
    "file_name" : file_name + current_time + strByte + strconv.Itoa(rand_number) + "." + format_file,
    "file_ext" : format_file,
    "created_by" : created_by.(string),
    "type" : u_type,
  })
}

func UploadTmp(file *multipart.FileHeader, c *gin.Context)string{
  // single file
  t := time.Now()
  if _, err := os.Stat("./public"); os.IsNotExist(err) {
    dErr := os.Mkdir("./public", 0755)
    if dErr != nil {
      gin.DefaultWriter.Write([]byte("ERROR : "+dErr.Error()))
    }
  }
  dest := "./public/tmp/"
  if _, err := os.Stat(dest); os.IsNotExist(err) {
    dErr := os.Mkdir(dest, 0755)
    if dErr != nil {
      gin.DefaultWriter.Write([]byte("ERROR : "+dErr.Error()))
    }
  }
  current_time := t.Format("20060102150405")
  rand_number := rand.Intn(9999)
  slice_data := strings.Split(file.Filename, ".")
  format_file := slice_data[len(slice_data) - 1]
  dest = dest + current_time + RandStringBytes(5) + strconv.Itoa(rand_number) + "." + format_file
  // Upload the file to specific dst.
  if err := c.SaveUploadedFile(file, dest); err == nil {
    return dest
  }
  return "error"
}

func UploadGet(id int)map[string]string{
  db := database.InitDB()
  var filename string
  row := db.QueryRow("SELECT file_name FROM uploads WHERE id = $1", id)
  row.Scan(&filename)
  res := map[string]string{
    "id":strconv.Itoa(id),
    "filename" : filename,
    "url" : os.Getenv("HOST") + "/public/uploads/" + filename,
  }
  return res
}

func UploadDelete(id int){
  db := database.InitDB()

  dbData := db.QueryRow("SELECT file_name FROM uploads WHERE id = $1", id)
  var (
    file_name string
  )
  dbData.Scan(&file_name)
  err := os.Remove("./public/uploads/"+file_name)
  if err != nil {
    gin.DefaultWriter.Write([]byte("ERROR : "+err.Error()))
  }
  if model.Delete("uploads", id) {
    gin.DefaultWriter.Write([]byte("ERROR : FAIL DELETE UPLOADS DATA"))
  }
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

