
# WebAPI（デスク環境の写真投稿サイト）

設計についてはこちらで解説しています
https://qiita.com/ryoh07/items/8ebac006c5294b9b3f58

## API

```
http://localhost:8080

ユーザの取得
GET /user/{id}

ユーザの新規作成
POST /user

ユーザの更新
POST /user/edit

ユーザの削除
DELETE /user/{id}

写真一覧の取得
GET /photo

写真の詳細取得
GET /photo/{id}

写真の新規作成
POST /photo

写真の更新
POST /photo/edit

自分が投稿した写真の取得
GET /photo/myphoto/{id}

自分がいいねした写真の取得
GET /photo/mylike/{id}

写真のいいねを付与
PUT /photo/{photoid}/{userid}/like

写真のいいねを削除
DELETE /photo/{photoid}/{userid}/like
```

# 使用例
## root（トップ画面）
トップ画面には投稿された写真の一覧が並びます。

![](https://i.imgur.com/rZlzNaZ.png)

#### エンドポイント
```
 GET /photo  
```

#### クエリパラメータ

| 項目 | 指定例       |備考       |
|:-------- |:---------- |:---------- |
| keyword   | ゲーマー    |投稿文章の部分検索
| tag_ids       | タグID　　| タグ検索(投稿ページ内のタグを選択)
| sort      | like  　　　|いいね順
| page       | 2　　|ページ番号

#### レスポンスボディ (JSON)
```json
{
    "photo_counts": 65,
    "photo_list": [
        {
            "id": "56c2646e",
            "photo": "xxxxxxxxxx",
            "likes": 3,
            "created_at": "2022-07-29 18:10:33",
            "user": {
                "id": "99b530ef",
                "name": "太郎",
                "icon": "xxxxxxxxxx",
                "self_introduction": "テスト"
            }
        },
        
        ～略～
        
        {
            "id": "45966d21",
            "photo": "",
            "likes": 0,
            "created_at": "2022-08-01 20:50:07",
            "user": {
                "id": "99b530ef",
                "name": "太郎",
                "icon": "xxxxxxxxxx",
                "self_introduction": "初投稿です"
            }
        }
    ]
}
```

## photo（投稿写真ページ） 
トップ画面より写真カードをクリックすると投稿写真ページに遷移します。

**機能**
・デスク環境の投稿写真の表示
・投稿コメントの表示
・写真内で使用されているアイテム(マウス、椅子などの品名)の表示
・登録されたタグの表示
・いいね数表示
・いいね付与
・いいね削除

![](https://i.imgur.com/Se6f0uc.png)



#### エンドポイント
```
 GET  /photo/{photoid}
```
##### レスポンスボディ (JSON)
```json
{
    "id": "b8ac3d4f",
    "photo": "xxxxxxxxxx",
    "contents": "テストです",
    "created_at": "2022-08-01 20:49:49",
    "likes": 10,
    "user": {
        "id": "99b530ef",
        "name": "太郎",
        "icon": "xxxxxxxxxx",
        "self_introduction": "よろしくおねがいします"
    },
    "item_counts": 2,
    "items": [
        {
            "id": 47,
            "name": "name1",
            "price": 1,
            "image": "xxxxxxxxxx"
        },
        {
            "id": 48,
            "name": "name2",
            "price": 2,
            "image": "xxxxxxxxxx"
        }
    ],
    "tags": [
        {
            "id": 3,
            "name": "ゲーマ"
        },
        {
            "id": 2,
            "name": "大学生"
        }
    ]
}
```

いいね付与  
```
 PUT  /photo/{photoid}/{userid}/like
```
いいね解除  
```
 DELETE  /photo/{photoid}/{userid}/like
```

## upload（アップロード画面） 
写真の投稿を行うページです。

**機能**
・写真と付随データ(投稿文、タグ、使用アイテム)を投稿する  

![](https://i.imgur.com/iWcQpEF.png)


#### エンドポイント
 ```
 POST /photo
 ```
 

#### リクエストボディ (JSON)
```json
{
    "photo": "xxxxxxxxxx",
    "contens": "test",
    "user_id" "99b530ef",
    "items":[
        {
            "name": "○○マウス",
            "price": 2000,
            "imege": "xxxxxxxxxx",
        },
        {
            "name": "○○キーボード",
            "price": 8000,
            "imege": "xxxxxxxxxx",
        }
    ],
    "tags":[
        {
            "name": "ゲーマ",
        },
        {
            "name": "大学生",
        }
    ]
}
```

## **upload/edit（編集画面）** 
投稿写真の編集画面

**機能**
・写真の付随データを編集、更新する(投稿文、タグ、使用アイテム)
 ![](https://i.imgur.com/0z1kWpz.png)
 
 #### エンドポイント
 ```
 POST /photo/edit
 ```
 #### リクエストボディ (JSON)
```json
{
    "photo_id": "b8ac3d4f",
    "contens": "更新しました",
    "user_id": "99b530ef",
    "items":[
        {
            "id" : 2,
            "name": "xxxxxxxxxx",
            "price": "xxxxxxx"
            "imege": "xxxxxxx",
        },
        {
            "id" : 0,  新しく作成したアイテムはidにゼロ値
            "name": "xxxxxxxxxx",
            "price": "xxxxxxx"
            "imege": "xxxxxxx",
        }
    ],
    "tags":[
        {
            "name": "ゲーマ",
        },
        {
            "name": "大学生",
        }
    ]
}
```

## **user（ユーザ画面）** 
**機能**
・アイコン、名前、自己紹介の表示
・投稿写真一覧の表示

![](https://i.imgur.com/mPA5OIi.png)


 #### エンドポイント
```
 GET /user/{id}
```  
#### レスポンスボディ (JSON)

```json
{
    "user_id": "99b530ef",
    "name": "太郎",
    "icon": "xxxxxxxx"
    "self_introduction": "自己紹介テスト",
    "photo_counts": "2",
    "photo_list": [
        {
            "id": "0c07725c",
            "photo": "xxxxxxxxxx",
            "created_at": "2022-08-01 20:25:15",
        },
        {
            "id": "367c34dc",
            "photo": "xxxxxxxxxx",
            "created_at": "2022-08-01 20:24:51",
        }
    ]
}
```

## **user/edit（ユーザ編集画面）** 
ユーザ情報の編集画面
![](https://i.imgur.com/bjOmgUC.png)

変更保存  
```
 POST /user/{id}/edit
```
#### リクエストボディ (JSON)

```json
{
    "user_id": "99b530ef",
    "name": "更新太郎",
    "icon": "xxxxxxxx"
    "self_introduction": "更新自己紹介テスト",
}
```


## **mypage（マイページ）** 
マイページ画面
![](https://i.imgur.com/4ApNHDl.png)



## **mypage/myphoto（投稿した写真）** 

自分が投稿した写真の一覧画面

**機能**
・写真一覧の取得
![](https://i.imgur.com/m5cC9FU.png)


 #### エンドポイント
```
 GET /photo/myphoto/{id}
```  
#### レスポンスボディ (JSON)
```json
"photo_counts": "2",
"photo_list": [
    {
        "id": "0c07725c",
        "photo": "xxxxxxxxxx",
        "created_at": "2022-08-01 20:25:15",
    },
    {
        "id": "367c34dc",
        "photo": "xxxxxxxxxx",
        "created_at": "2022-08-01 20:24:51",
    }
]
```

## **mypage/mylike（いいねした写真）** 

**機能**
・自分がいいねした投稿の一覧の表示  
  
![](https://i.imgur.com/bkJ0AiE.png)
 #### エンドポイント
```
 GET /photo/{id}/mylike
```  
#### レスポンスボディ (JSON)
```json
"photo_counts": "1800",
"photo_list": [
    {
        "id": "b590ddab",
        "photo": "xxxxxxxxxx",
        "created_at": "2022-08-01 20:17:57",
        "user": { 
            "id": "c736a1c1",
            "name": "テスト",
            "icon": "xxxxxxxx",
        }
    },
    {
        "id": "c4de96e2",
        "photo": "xxxxxxxxxx",
        "created_at": "2022-08-01 20:17:48",
        "user": { 
            "id": "c736a1c1",
            "name": "テスト",
            "icon": "xxxxxxxx",
        }
    }
]
```

