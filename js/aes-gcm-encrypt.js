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
  async function encryptMessage(key, iv) {
    let encoded = getMessageEncoding();

    let ciphertext = await window.crypto.subtle.encrypt({
        name: "AES-GCM",
        iv: iv
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
        document.querySelector(".link").innerHTML = "Link: <a href='" + url + "'>link</a>";
      }
    }
    let _iv = document.querySelector(".iv").innerText;
    let _key = document.querySelector(".encryption-key").innerText;
    let _eol = document.querySelector("#eol option:checked").value;
    xhr.send("iv="+_iv+"&key="+_key+"&content="+ctBase64+"&eol="+_eol);
  }

  /*
  Export the given key and write it into the "exported-key" space.
  */
  async function exportCryptoKey(key) {
    const exported = await window.crypto.subtle.exportKey(
      "jwk",
      key
    );
    const exportKeyOutput = document.querySelector(".encryption-key");
    exportKeyOutput.textContent = window.btoa(JSON.stringify(exported, null, " "));
    console.log(JSON.stringify(exported, null, " "));
  }

  /*
  Generate an encryption key, then set up event listeners
  on the "Encrypt" and "Decrypt" buttons.
  */
  window.crypto.subtle.generateKey({
      name: "AES-GCM",
      length: 256,
    },
    true,
    ["encrypt", "decrypt"]
  ).then((key) => {
    const encryptButton = document.querySelector(".encrypt-button");
    exportCryptoKey(key);
    // The iv must never be reused with a given key.
    let iv = window.crypto.getRandomValues(new Uint8Array(12));
    iv = window.crypto.getRandomValues(new Uint8Array(12));
    var htmlIv = document.querySelector(".iv");
    htmlIv.innerText = window.btoa(iv.toString());
    encryptButton.addEventListener("click", () => {
      encryptMessage(key, iv);
    });
    // Just for fun... Encrypt the message in case of update
    // Since encryptMessage also send data... In case of use we have to create
    // a dedicated function for sending message.
    //document.querySelector('textarea').addEventListener('input', function (event) {
    //  if (event.target.value.length != 0) {
    //    encryptMessage(key, iv)
    //  }
    //});

  });

})();