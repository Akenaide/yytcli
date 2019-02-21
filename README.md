# yytcli

yytcli fetch WS cards informations from yuyutei website.

You can use it to create a snapshot of today price with:

    yytcli getcards
    
  or
  
    yytcli getprice
          
  to get a playset price

Usage:
  yytcli [command]

Available Commands:

  getcards    get card infos from yyt
  
  getprice    get playset price (RR, R, U, C, CC, CR)
  
  help        Help about any command

Flags:

  --config string        config file (default is $HOME/.ycli.yaml)
      
  -h, --help                 help for yytcli
  
  -k, --kizu                 Get only damaged informations
  
  -s, --series stringArray   Default fetch all series
  
  -t, --toggle               Help message for toggle
  

Use "yytcli [command] --help" for more information about a command.
