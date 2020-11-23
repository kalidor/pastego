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

        <!--link href="/assets/stylesheets/bootstrap.min.css" rel="stylesheet">
        <link href="/assets/stylesheets/sticky-footer.css" rel="stylesheet"-->

    </head>
    <body>
        <div class="container">
            <div class="page-header">
                <a class="h1" href="/">PasteGo</a>
            </div>
            
            
            <p>New Paste:</p>
    
     

<form action="/create" method="POST" >
    
        <textarea name="content" class="form-control" rows="20"></textarea>
        <br />
        <button type="submit" class="btn btn-default pull-right">Create</button>
        <div class="col-sm-2 pull-right">
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
        </div>
</form>
        </div>
        <div id="footer">
            <div class="container">
                <p class="text-muted">Developed by gch - 2020</p>
            </div>
        </div>
    </body>
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

        <!--link href="/assets/stylesheets/bootstrap.min.css" rel="stylesheet">
        <link href="/assets/stylesheets/sticky-footer.css" rel="stylesheet"-->

    </head>
    <body>
        <div class="container">
            <div class="page-header">
                <a class="h1" href="/">PasteGo</a>
            </div>
            
            
            <p>{{ .Pasteid }}</p>
    
     

        <textarea name="content" class="form-control" rows="20">{{ .Content}}</textarea>
        <br />

        </div>

        <div id="footer">
            <div class="container">
                <p class="text-muted">Developed by gch - 2020</p>
            </div>
        </div>
    </body>
</html>
`
