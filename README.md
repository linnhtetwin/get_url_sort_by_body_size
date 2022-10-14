# Visit Url and Write URL : Resp Body Size Pairs in text file (sorted by small to large body size)

## To Test
```
go run urlcheck.go
```
Step 1,
Please locate your text file path where you added List of Urls that you want to check!
For example : /home/url.txt

Step 2,
Go and Check the same text file

You don't need to worry about the prefix,
if your url don't have http prefix, it will automatically add.
if your url have http prefix, it will not.

This version of Cli tool replace the same txt file to avoid the filepath error on different os.
So, It will visit each url and rewrite the output result in your same text file.

## Output Format

https://.......     Resp Body Size : 12.....

https://.......     Resp Body Size : 12.....

https://.......     Resp Body Size : 12.....

## To install and run in cli

```
go install urlcheck.go
urlcheck
```

## Possible Errors
1. Your url have extra space at the end of each line
2. Non-url texts in your file