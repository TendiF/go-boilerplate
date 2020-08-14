package model

import (
  "../database"
  "fmt"
  "strconv"
)
func Add(tabel string,dataMap map[string]string) int{
  q := "INSERT INTO public."+tabel+" "
  i := 0
  qv := " VALUES "
  values := []interface{}{}
  for k, v := range dataMap {
    i++
    if i == 1 {
      q += "("
      qv += "("
    }
    q += k
    // qv += "'"+v+"'"
    qv += "$"+strconv.Itoa(i)
    if len(dataMap) == i {
      q += ")"
      qv += ")"
      q += qv
    } else {
      q += ", "
      qv += ", "
    }
    values = append(values, v)
  }
  q += " RETURNING id"
  fmt.Println(q)
  db := database.InitDB()
  lastId := 0
  err := db.QueryRow(q,values... ).Scan(&lastId)
  if err != nil {
    return 0
  }
  return lastId
}

func Delete(tabel string, id int)bool{
  q := "DELETE FROM "+tabel+" WHERE id = $1"
  db := database.InitDB()
  _, err  := db.Exec(q,id)
  if err != nil {
    return false
  } else {
    return true
  }
}

func DeleteSt(tabel string, id int) bool {
  db := database.InitDB()
  q := "UPDATE "+tabel+" SET status = -2 WHERE id = $1"
  _, err := db.Exec(q, id)
  if err != nil {
    return false
  } else {
    return true
  }
}

func Update(tabel string, dataMap map[string]string, id int) bool {
  db := database.InitDB()
  q := "UPDATE "+tabel
  // SET column1 = value1, column2 = value2
  // WHERE condition;"
  i := 0;
  params := []interface{}{}
  for k, v := range dataMap {
    i++;
    if i == 1 {
      q += " SET "
    }
    q += k+" = $"+strconv.Itoa(i)
    if len(dataMap) != i {
      q += ","
    }
    params = append(params, v)
  }
  q += " WHERE id = $"+strconv.Itoa(i+1)
  params = append(params, id)
  _, err := db.Exec(q, params...)
  if err != nil {
    return false
  } else {
    return true
  }
}

func Get(statement map[string]string)([]interface{}, map [string]int){
  db := database.InitDB()
  q := ""
  if statement["select"] != "" {
    q += "SELECT "+statement["select"]+" "
  } else {
    q += "SELECT * "
  }
  q += "FROM "+statement["tabel"]
  if statement["join"] != "" {
    q += " JOIN " + statement["join"]
  }
  if statement["left_join"] != "" {
    q += " LEFT JOIN " + statement["left_join"]
  }
  if statement["where"] != "" {
    q += " WHERE " + statement["where"]
  }
  if statement["group_by"] != "" {
    q += " GROUP BY " + statement["group_by"]
  }
  if statement["order_by"] != "" {
    q += " ORDER BY " + statement["order_by"]
  }
  if statement["query"] != "" {
    q = statement["query"]
  }
  if statement["limit"] != "" || statement["offset"] != "" {
    if (statement["limit"] == "") || (statement["limit"] == "0") {
      statement["limit"] = "5"
    }
    if (statement["offset"] == "") || (statement["offset"] == "0") {
      statement["offset"] = "0"
    }
    offset, _ := strconv.Atoi(statement["offset"])
    limit, _ := strconv.Atoi(statement["limit"])

    q += " OFFSET $1 LIMIT $2"
    rows, err := db.Query(q, statement["offset"], statement["limit"]) // Note: Ignoring errors for brevity
    if err != nil {
      panic(err)
    }
    cols, _ := rows.Columns()
    data := []interface{}{}
    defer rows.Close()
    for rows.Next() {
      columns := make([]interface{}, len(cols))
      columnPointers := make([]interface{}, len(cols))
      for i, _ := range columns {
          columnPointers[i] = &columns[i]
      }

      rows.Scan(columnPointers...);

      m := make(map[string]interface{})
      for i, colName := range cols {
          val := columnPointers[i].(*interface{})
          m[colName] = *val
      }

      data = append(data, m)
    }

    rows, err = db.Query(q, offset, limit+1)
    count := 0
    for rows.Next() {
      count++
    }

    is_last := 0

    if count != limit+1 {
      is_last = 1
    }
    pagination := map [string]int{
      "current_page" : (offset / limit) + 1,
      "current_data" : len(data),
      "offset" : offset,
      "limit" : limit,
      "is_last" : is_last,
    }
    return data, pagination
  }else {
    rows, err := db.Query(q)
    if err != nil {
      panic(err)
    }
    cols, _ := rows.Columns()
    data := []interface{}{}
    defer rows.Close()
    for rows.Next() {
      columns := make([]interface{}, len(cols))
      columnPointers := make([]interface{}, len(cols))
      for i, _ := range columns {
          columnPointers[i] = &columns[i]
      }

      rows.Scan(columnPointers...);

      m := make(map[string]interface{})
      for i, colName := range cols {
        val := columnPointers[i].(*interface{})
        if colName == "rate" {
          data := *val
          m[colName] = string(data.([]uint8))
        } else {
          m[colName] = *val
        }
      }

      data = append(data, m)
    }

    is_last := 1
    pagination := map [string]int{
      "current_page" : 1,
      "current_data" : len(data),
      "offset" : 0,
      "limit" : len(data),
      "is_last" : is_last,
    }
    return data, pagination
  }
}

func Exist(statement map[string]string)bool{
  res, _ := Get(statement)
  if len(res) == 0 {
    return false
  }
  return true
}
