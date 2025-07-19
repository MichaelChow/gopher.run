---
title: "5.8 gorm"
date: 2025-05-26T23:38:00Z
draft: false
weight: 5008
---

# 5.8 gorm

- 标准库写原生sql：[https://go.dev/doc/tutorial/database-access](https://go.dev/doc/tutorial/database-access)
    ```go
    // albumsByArtist queries for albums that have the specified artist name.
    func albumsByArtist(name string) ([]Album, error) {
        // An albums slice to hold data from returned rows.
        var albums []Album
        rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
        if err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        defer rows.Close()
        // Loop through rows, using Scan to assign column data to struct fields.
        for rows.Next() {
            var alb Album
            if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
                return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
            }
            albums = append(albums, alb)
        }
        if err := rows.Err(); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        return albums, nil
    }
    ```
- 第三方ORM库：GORM、XORM 、ent.
- GORM官方文档：[https://gorm.io/zh_CN/docs/create.html](https://gorm.io/zh_CN/docs/create.html) 跟着教程写一遍CRUD和autoMigrate？


