package main

const INDEX = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="description" content="Very Simple Pastebin">
        <meta name="author" content="gch">
        <link rel="shortcut icon" href="/assets/images/favicon.png">

        <title>PasteGo - Create paste</title>

        <link href="/css/bootstrap.min.css" rel="stylesheet">

    </head>
    <body>
        <div class="container">
            <div class="page-header">
                <a class="h1" href="/">PasteGo</a>
            </div>
            <p>New Paste:</p>
    
        <textarea name="content" class="form-control" rows="20"></textarea>
        <br />
        <div class="col-sm-5 pull-right">
            <select name="eol" class="form-control" id="eol">
                <option value="10">10 min</option>
                <option value="30">30 min</option>
                <option value="60">01 h</option>
                <option value="720">12 h</option>
                <option selected value="1440">01 j</option>
                <option value="2880">02 j</option>
                <option value="10080">07 j</option>
                <option value="21600">15 j</option>
                <option value="43200">30 j</option>
            </select>
        Password: <input type="text" class="password" placeholder="password" />
        <input class="encrypt-button" type="button" value="Encrypt" />
        <span class="link"></span>
        </div>
        <input type="hidden" class="iv" />
        <input type="hidden" class="ciphertext" />
        </div>
        <div id="footer">
            <div class="container">
                <p class="text-muted">Developed by gch - 2020</p>
            </div>
        </div>
    </body>
    <script src="js/aes-gcm-encrypt.js"></script>
</html>`

const VIEW = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="description" content="Very Simple Pastebin">
        <meta name="author" content="gch">
        <link rel="shortcut icon" href="/assets/images/favicon.png">

        <title>PasteGo - Create paste</title>

        <link href="/css/bootstrap.min.css" rel="stylesheet">
        <link href="/css/sticky-footer.css" rel="stylesheet">

    </head>
    <body>
        <div class="container">
            <div class="page-header">
                <a class="h1" href="/">PasteGo</a>
            </div>
            
                <div class="row">
        <div class="col-md-4"><strong>Posted at </strong>{{ .TimeStart }}</div>
        <div class="col-md-4"><strong>Ends at </strong>{{ .TimeStop }}</div>
        <div class="col-md-4"><strong>Raw </strong><a href="http://{{ .Url }}/raw/{{ .Pasteid }}">{{ .Pasteid }}</a></div>
    </div>
    <br />
        <textarea name="content" class="form-control" id="message" rows="20">{{ .Content }}</textarea>
        Password: <input type="text" class="password" placeholder="password"/>
        <input class="decrypt-button" type="button" value="Decrypt" />
        </div>
        <input type="hidden" class="iv" value="{{ .Iv }}"/>
        <div id="footer">
            <div class="container">
                <p class="text-muted">Developed by gch - 2020</p>
            </div>
        </div>
    </body>
    <script src="/js/aes-gcm-decrypt.js"></script>
</html>
`
