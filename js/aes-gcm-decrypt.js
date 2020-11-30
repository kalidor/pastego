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
  async function decryptMessage(mkey) {
    var htmlIv = window.atob(document.querySelector(".iv").value);
    var myArray = htmlIv.split(",");
    for (var i = 0; i < myArray.length; i++) {
      myArray[i] = parseInt(myArray[i], 10);
    }
    const s = new Set(myArray);
    let _iv = new Uint8Array(12);
    _iv = Uint8Array.from(s);

    const messageBox = document.querySelector(".ciphertext-value");
    const ctStr = window.atob(messageBox.value).match(/[\s\S]/g);
    const ctUint8 = new Uint8Array(ctStr.map((char) => char.charCodeAt(0)));
    const decrypted = await window.crypto.subtle.decrypt({
        name: "AES-GCM",
        iv: _iv
      },
      mkey,
      ctUint8
    );
    await new Promise(r => setTimeout(r, 1000));
      let dec = new TextDecoder();
      const decryptedValue = document.querySelector("textarea");
      decryptedValue.textContent = dec.decode(decrypted);
  }

  function importSecretKey() {
    var htmlKey = window.atob(document.querySelector(".encryption-key").value);
    window.crypto.subtle.importKey(
        "jwk",
        JSON.parse(htmlKey),
        "AES-GCM",
        true,
        ["encrypt", "decrypt"]
      )
      .then((mkey) => {
        decryptMessage(mkey);
      })
      .catch(function (err) {
        console.log(err);
      })
  };
  importSecretKey();

})();