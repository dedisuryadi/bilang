[![Go](https://github.com/dedisuryadi/bilang/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/dedisuryadi/bilang/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedisuryadi/bilang.svg)](https://pkg.go.dev/github.com/dedisuryadi/bilang)

# Bilang
Bilang (Bahasa is not a language) adalah sebuah interpreter untuk bahasa pemrograman berbahasa indonesia.

**Bilang** tidak mengikuti kaidah berbahasa Indonesia yang baik dan benar ketika melakukan terjemahan dari Bahasa Inggris ke Bahasa Indonesia, 
hal tersebut disengaja agar tidak kehilangan arti aslinya, kita lebih paham World Wide Web (www) daripada Waring Wera Wanua.


Fitur yang sudah bisa digunakan:
- [x] REPL
```
    go run .
```
atau 
```
    cat script.bi | go run .
```

- [x] Variabel
```
var a = "halo dunia"
println(a)

```

- [x] Konstanta   
```
konst a = "halo dunia"
println(a)
```

- [x] Fungsi adalah *First class citizen*
```
var reduce = fn(arr, init, f) {
    var iter = fn(arr, hasil) {
        jika (panjang(arr) == 0) { pilih hasil };
        iter(ekor(arr), f(hasil, awal(arr)))
    }
    iter(arr, init)
}

var sum = arr => reduce(arr, 0, fn(init, nilai){ init+nilai })	
sum([1,2,3,4,5])
```


- [x] Pipe operator
```
var reduce = fn(arr, init, f) {
    var iter = fn(arr, hasil) {
        jika (panjang(arr) == 0) { pilih hasil };
        iter(ekor(arr), f(hasil, arr |> awal))
    }
    iter(arr, init)
}

var sum = arr => reduce(arr, 0, fn(init, nilai){ init+nilai })
var map = fn(arr, f) {
    var iter = fn(arr, akum) {
        jika (panjang(arr) == 0) { pilih akum }
        var hasil = push(akum, arr |> awal |> f)
        iter(arr |> ekor, hasil)
    }
    iter(arr, [])
}

var a = [1,2,3,4,5]
var ganda = x => x*2
map(a, ganda) |> sum

```

- [x] Strict typing
```shell

$ echo '
 var a = "string"
 a = 10
 ' | go run .
ERROR: perubahan tipe variabel a dari STRING menjadi FLOAT tidak diizinkan

```
    
- [X] Switch statement
```
var rgb_ke_hsl = fn(arr) {
    var r = arr[0]/255
    var g = arr[1]/255
    var b = arr[2]/255
    var max = math.Max(r, math.Max(g, b))
    var min = math.Min(r, math.Min(g, b))
    var h = 0
    var l = (max+min)/2
    var d = max - min
    var s = jika (l > 1/2) {
        d / (2 - max - min)
    } atau {
        d / (max + min)
    }

    h = pilah max {
        r -> (g - b) / d + jika (g < b) { 6 } atau { 0 }
        g -> (b - r) / d + 2
        b -> (r - g) / d + 4
    }

    pilih [h, s, l]
}

var merah_rgb = [255, 0, 0]
var merah_hsl = rgb_ke_hsl(merah_rgb)

stdout("merah_rgb = ", merah_rgb, "\n")
stdout("merah_hsl = ", merah_hsl, "\n")

```

- [x] Dan lainnya


TODO:
- [ ] Komentar
- [ ] Error handling
- [ ] Tipe Data Integer
- [ ] Notasi angka float (dan eksponen)
- [ ] Standard library
- [ ] Testing ala go test
- [ ] Notasi pendek variabel menggunakan `:=` seperti Go
- [ ] Modul sistem (ekspor & impor)


Repo ini pada awalnya dibuat sebagai tempat latihan saat membaca buku **Writing An Interpreter In Go** 
maka dari itu akan ada banyak kesamaan struktur kode dengan **Monkey Language**. Pull request welcome. 
