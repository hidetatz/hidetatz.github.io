about commandline tool to generate slide from markdown---2018-07-24 19:09:52

We developers sometimes need to create slides for presentation, talking at conferences.
Always, how do you create it? There are several ways, [marp](https://yhatt.github.io/marp/), [remarkjs](https://remarkjs.com/#1), and so on.
If you used to use markdown and commandline, you can try [slidegen](https://github.com/yagi5/slidegen). Yes, I am an author of it.

---  
  
### Simple and minimal command  
  
![](https://raw.githubusercontent.com/ygnmhdtt/slidegen/master/samples/demo.gif)
  
All you have to do is specifying file name just like this.  
  
```sh  
$ slidegen your/markdown/file.md
```  
  
Then, you can get output.pdf at your current directory.

---

### Markdown format

All markdown syntaxes are supported that are introduced [here](https://guides.github.com/features/mastering-markdown/).
Delimiter of each page is `---` .

---

### PDF style

The simplest and easiest-to-read markdown style, GFM(github-flavored-markdown) will be applied on your PDF.
Always I need only it, but if you want to use another css, please fork and customize it.

---

### Import from Gist

You can import from Gist URL. Just give `-g` option like

```sh
$ slidegen -g https://gist.github.com/
```

Of course, gist must contain `---` for each page delimiter

---

### Please try

I want you to use it and, if you like, please star it.
If any problems or feedback, please create a [issue](https://github.com/ygnmhdtt/slidegen/issues).
Pull requests are always welcome!

---

### Finally

I created slide from this entry. (PDF cannot be posted on medium, I converted into images.)

![](https://raw.githubusercontent.com/yagi5/slidegen/master/medium/1.png)
![](https://raw.githubusercontent.com/yagi5/slidegen/master/medium/2.png)
![](https://raw.githubusercontent.com/yagi5/slidegen/master/medium/3.png)
![](https://raw.githubusercontent.com/yagi5/slidegen/master/medium/4.png)
![](https://raw.githubusercontent.com/yagi5/slidegen/master/medium/5.png)
