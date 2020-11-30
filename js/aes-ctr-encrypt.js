(() => {

  /*
  Fetch the contents of the "message" textbox, and encode it
  in a form we can use for the encrypt operation.
  */
  function getMessageEncoding() {
    const messageBox = document.querySelector("textarea");
    let message = messageBox.value;
    let enc = new TextEncoder();
    return enc.encode(message); // encode plaintext as UTF-8
  }

  /* Encrypt message */
  async function encryptMessage() {
    let encoded = getMessageEncoding();
    // The iv must never be reused with a given key.
    const iv = window.crypto.getRandomValues(new Uint8Array(16));
    var htmlIv = document.querySelector(".iv");
    htmlIv.innerText = window.btoa(iv.toString());

    const pwUtf8 = new TextEncoder().encode(document.querySelector(".password").innerTxt);  // encode password as UTF-8
    const pwHash = await crypto.subtle.digest('SHA-256', pwUtf8);                           // hash the password
    const alg = { name: 'AES-CTR'};
    const key = await crypto.subtle.importKey('raw', pwHash, alg, false, ['encrypt']);      // use pw to generate key

    let ciphertext = await window.crypto.subtle.encrypt({
        name: 'AES-CTR',
        counter: iv,
        length: 64
      },
      key,
      encoded
    );
    const ciphertextValue = document.querySelector(".ciphertext");
    ciphertextValue.classList.add('fade-in');
    ciphertextValue.addEventListener('animationend', () => {
      ciphertextValue.classList.remove('fade-in');
    });
    const ctArray = Array.from(new Uint8Array(ciphertext)); // ciphertext as byte array
    ctStr = ctArray.map(byte => String.fromCharCode(byte)).join(''); // ciphertext as string
    const ctBase64 = window.btoa(ctStr); // encode ciphertext as base64
    ciphertextValue.textContent = ctBase64;

    // Send data
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/create", true);
    //Send the proper header information along with the request
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    xhr.onreadystatechange = function() { // Call a function when the state changes.
      if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
        var url = xhr.responseURL;
        document.querySelector(".link").innerHTML = "Link: <a href='" + url + "'>"+ url +"</a>";
      }
    }
    let _iv = document.querySelector(".iv").innerText;
    let _eol = document.querySelector("#eol option:checked").value;
    xhr.send("iv="+_iv+"&content="+ctBase64+"&eol="+_eol);
  }
  const encryptButton = document.querySelector(".encrypt-button");
    encryptButton.addEventListener("click", () => {
    encryptMessage();
  });

})();
