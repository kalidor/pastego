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

  /*
  Fetch the ciphertext and decrypt it.
  Write the decrypted message into the "Decrypted" box.
  */
  async function decryptMessage() {
    var htmlIv = window.atob(document.querySelector(".iv").value);
    var myArray = htmlIv.split(",");
    for (var i = 0; i < myArray.length; i++) {
      myArray[i] = parseInt(myArray[i], 10);
    }
    const s = new Set(myArray);
    let _iv = new Uint8Array(16);
    _iv = Uint8Array.from(s);

    const messageBox = document.querySelector("textarea");
    const ctStr = window.atob(messageBox.value).match(/[\s\S]/g);
    const ctUint8 = new Uint8Array(ctStr.map((char) => char.charCodeAt(0)));

    const pwUtf8 = new TextEncoder().encode(document.querySelector(".password").innerTxt);  // encode password as UTF-8
    const pwHash = await crypto.subtle.digest('SHA-256', pwUtf8);                           // hash the password
    const alg = { name: 'AES-CTR'};
    const key = await crypto.subtle.importKey('raw', pwHash, alg, false, ['decrypt']);      // use pw to generate key

    const decrypted = await window.crypto.subtle.decrypt({
        name: "AES-CTR",
        counter: _iv,
        length: 64
      },
      key,
      ctUint8
    );
    //await new Promise(r => setTimeout(r, 1000));
    const clean = new TextDecoder().decode(decrypted);
    const decryptedValue = document.querySelector("textarea");
    decryptedValue.textContent = clean;
  }
  const decryptButton = document.querySelector(".decrypt-button");
  decryptButton.addEventListener("click", () => {
    decryptMessage();
  });
  

})();