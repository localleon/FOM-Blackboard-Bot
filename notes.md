# Steps for Login:
- Get Login JSESSIONID
- Perform Login on Login.do with Username and Shit (you get a session cookie)
- Append session cookie and then perform get requests

## Login Context Form Return

<!DOCTYPE html>
<html lang="de">
<head>
    <title>Online-Campus Login</title>
    
    <meta charset="utf-8">
    
	
    <script>
        function init(){
            var placeholderSupported = !!( 'placeholder' in document.createElement('input') );
            if(!placeholderSupported){
                document.body.className += ' showLabels';
            }
            document.login.name.focus();
        }
    </script>
	
	<link rel="stylesheet" type="text/css" href="/nfcampus/css/login.css">
	
</head>

<body onload="init()" class="fom">

	
	
    <form method="post" action="/nfcampus/Login.do" accept-charset="UTF-8" name="login">
    	<input type="hidden" name="crt" value="19453">
    	<input type="hidden" name="assl" value="">
        <input name="iehack" type="hidden" value="â ">
        <input name="quelle" type="hidden" value="LoginForm-FOM">
        <input name="i" type="hidden" value="fom">
        <strong>Online-Campus</strong>
        <div id="inForm">
        	            	
            <div id="cell_name">
            
                <label for="name">Benutzername</label>
                <input type="text" class="text" placeholder="Benutzername" id="name" name="name" />
            </div>
            <div id="cell_pass">
                <label for="pass">Passwort</label>
                <input type="password" class="text" placeholder="Passwort" id="pass" name="password" />
            </div>
            <div id="cell_submit">
                <a href="/nfcampus/reset.jsp?i=fom">Passwort vergessen?</a><input type="submit" value="&raquo; Login" />
            </div>
            <div id="cell_forgot">
                <a href="/nfcampus/pages/security.jsp?omitLink=true" target="hinweis" onclick="window.open('', 'hinweis', 'location=0,status=0,scrollbars=0,width=550,height=270')">Hinweise zur Sicherheit</a>
            </div>
        </div>
    </form>
</body>
</html>