var CACHE_NAME = "pwa-example-v1";

console.log("loading: serviceworker");

var urlsToCache = [
  ".",
  "./offline.html",
  "./wasm_exec.js",
  "./main.wasm",
  "./serviceworker.js",
  "./assets/favicon.svg",
  "./assets/favicon.png",
  "https://unpkg.com/spectre.css/dist/spectre.min.css",
  "https://unpkg.com/spectre.css/dist/spectre-exp.min.css",
  "https://unpkg.com/spectre.css/dist/spectre-icons.min.css",
];

// 残したいキャッシュのバージョンをこの配列に入れる
// 基本的に現行の1つだけでよい。他は削除される。
const CACHE_KEYS = [CACHE_NAME];

self.addEventListener("install", function (event) {
  console.log("install: serviceworker");
  event.waitUntil(
    caches
      .open(CACHE_NAME) // 上記で指定しているキャッシュ名
      .then(function (cache) {
        return cache.addAll(urlsToCache); // 指定したリソースをキャッシュへ追加
        // 1つでも失敗したらService Workerのインストールはスキップされる
      })
  );
});

//新しいバージョンのServiceWorkerが有効化されたとき
self.addEventListener("activate", (event) => {
  console.log("activate: serviceworker");
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys
          .filter((key) => {
            return !CACHE_KEYS.includes(key);
          })
          .map((key) => {
            // 不要なキャッシュを削除
            return caches.delete(key);
          })
      );
    })
  );
});

self.addEventListener("fetch", function (event) {
  var online = navigator.onLine;

  // ファイルパス ~/test.htmlにアクセスすると、このファイル自体は無いがServiceWorkerがResponseを作成して表示してくれる
  if (event.request.url.indexOf("test.html") != -1) {
    return event.respondWith(
      new Response("任意のURLの内容をここで自由に返却できる")
    );
  }

  if (online) {
    console.log("ONLINE");
    //このパターンの処理では、Responseだけクローンすれば問題ない
    //cloneEventRequest = event.request.clone();
    event.respondWith(
      caches.match(event.request).then(function (response) {
        if (response) {
          //オンラインでもローカルにキャッシュでリソースがあればそれを返す
          //ここを無効にすればオンラインのときは常にオンラインリソースを取りに行き、その最新版をキャッシュにPUTする
          return response;
        }
        //request streem 1
        return fetch(event.request)
          .then(function (response) {
            //ローカルキャッシュになかったからネットワークから落とす
            //ネットワークから落とせてればここでリソースが返される

            // Responseはストリームなのでキャッシュで使用してしまうと、ブラウザの表示で不具合が起こる(っぽい)ので、複製したものを使う
            cloneResponse = response.clone();

            if (response) {
              if (response || response.status == 200) {
                console.log("正常にリソースを取得");
                caches.open(CACHE_NAME).then(function (cache) {
                  console.log("キャッシュへ保存");
                  //初回表示でエラー起きているが致命的でないので保留
                  cache.put(event.request, cloneResponse).then(function () {
                    console.log("保存完了");
                  });
                });
              } else {
                return event.respondWith(
                  new Response("200以外のエラーをハンドリングしたりできる")
                );
              }
              return response;
            }
          })
          .catch(function (error) {
            return console.log(error);
          });
      })
    );
  } else {
    console.log("OFFLINE");
    event.respondWith(
      caches.match(event.request).then(function (response) {
        // キャッシュがあったのでそのレスポンスを返す
        if (response) {
          return response;
        }
        //オフラインでキャッシュもなかったパターン
        return caches.match("offline.html").then(function (responseNodata) {
          return responseNodata;
        });
      })
    );
  }
});
