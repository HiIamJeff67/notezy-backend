## 對於 util 的單元測試需求

### 架構

- 在位於 test/unit/util 底下
- 建立多個 test_some_util.go 每個對應到 app/util/ 內的一個 go 檔案
- 將所有測試資料放到 test/unit/util/testdata/ 底下
- 建立一個通用的 TestCase (已建立在 shared/type/test_case.go) :

```Go
// notezy-backend/shared/types/test_case.go
type TestCase [ArgType any, ReturnType any]struct {
  Args ArgType
  Returns ReturnType
}
```

### 實作步驟

- 首先先在 test/unit/util/ 底下先撰寫一個對應 app/util/ 底下的 Go 檔案
  - 假設現在在 app/util/some_util_function.go，那我們就要建立一個 test/unit/util/some_util_function_test.go，並在此撰寫測試相關的程式碼，其中該檔案包含 :
    - 以 Util 的 Function 名稱 + ArgType 作為被測試函式的參數型態，依照上面的範例就是 SomeUtilFunctionArgType
    - 以 Util 的 Function 名稱 + ReturnType 作為被測試函式的參數型態，依照上面的範例就是 SomeUtilFunctionReturnType
    - 以 Util 的 Function 名稱 + TestCase 作為被測試函式的 Input/Output 資料架構，依照上面的範例就是 SomeUtilFunctionTestCase，且使用剛剛建立的 SomeUtilFunctionArgType 以及 SomeUtilFunctionReturnType 指定他的 ArgStruct 以及 ReturnType，做出來的結果類似底下程式碼範例 :

```Go
// 範例中展示假設一個 SomeFunction 接收兩個 int 參數並回傳一個 int 回傳值。
type SomeUtilFunctionArgType = struct {
  A int
  B int
}
type SomeUtilFunctionReturnType = int
type SomeUtilFunctionTestCase  = TestCase[SomeUtilFunctionArgs, SomeUtilFunctionReturnType]
```

- 接著在 test/unit/testdata/ 底下建立關於這份 Go 檔案的 testdata 檔案，並將其命名為 some_util_function_testdata.json，其中的格式需要符合我們剛剛定義的 SomeUtilFunctionTestCase，舉例依照我們剛剛建立的 SomeUtilFunctionTestCase :

```json
[
  { "Args": { "A": 1, "B": 2 }, "Returns": 3 },
  { "Args": { "A": 5, "B": 7 }, "Returns": 12 }
]
```

- 然後繼續回到我們的 test/unit/util/some_util_function_test.go 內去撰寫 透過讀取我們剛剛寫在 testdata/ 內的 some_util_function_testdata.json 的資料來運行實際關於該 util 的測試程式碼。
- 注意 : 一個 util/some_util_function.go 裡面可能會有多個函數，請都針對他們重複上述針對 some_util_function 的流程，並確保測試這些函式的程式碼都被寫在同一個 some_util_function_test.go 檔案內，假設 some_util_function.go 裡面有 f1(), f2(), f3()，那麼你就要建立如下的結構 :

```MarkDown
- test/
  - unit/
    - testdata/
      - some_util_function_testdata/
        - f1_testdata.json
        - f2_testdata.json
        - f3_testdata.json
```

- 請依照上面要求，完成對應所有 app/util/ 的所有 Go 檔案的測試項目
