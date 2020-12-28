const version = "v1";

class Transaction {
	constructor() {
        this.txobj = new Object();
		this.txobj.type = 0;
        this.txobj.hash = version; 
		
		this.txobj.public_key = version;
		this.txobj.sig = version;

		this.txobj.block_hash = version;
		this.txobj.merkle_path = version;

		this.txobj.timestamp = Date.now();
		
		this.txobj.fields = new Array();
	}

	setType(t) {
		this.txobj.type = t;
	}

	setPublicKey(pk) {
		this.txobj.public_key = pk;
	}
	
	setTimestamp(tm) {
		this.txobj.timestamp = tm;
	}

	setSig(sig) {
		this.txobj.sig = sig;
	}
	
	addField(n, value) {
		//console.log("addOut:", value, bytecode);
		var field = new Object();
		field.n = n;
		//field.data = window.btoa(value);
		field.value = value;
		//console.log(out1.data);
		this.txobj.fields.push(field);
	}

	getBytes() {
		var txBytes = new Array();		
		var te = new TextEncoder("utf-8");
		
		var bType = this.txobj.type;
		txBytes.push(bType); // просто один байт
		// Пример перевода в байты числа
		//console.log(new Uint8Array(new Int32Array([1591512893]).buffer))
		
		var bPk = te.encode(this.txobj.public_key);
		txBytes = txBytes.concat([].slice.call(bPk));

		//var bTimestamp = te.encode(this.txobj.timestamp);
		var bTimestamp = new Uint8Array(new Int32Array([this.txobj.timestamp]).buffer);
		var bT0 = new Uint8Array(new Int32Array([0]).buffer);
		txBytes = txBytes.concat([].slice.call(bTimestamp));
		txBytes = txBytes.concat([].slice.call(bT0)); // TODO объединить в один 8-байтовый

		for(this.key in this.txobj.fields) {
			var bfn = new Uint8Array(new Int32Array([this.txobj.fields[this.key].n]).buffer);
			txBytes = txBytes.concat([].slice.call(bfn));
					
			var bfv = te.encode(this.txobj.fields[this.key].value);
			txBytes = txBytes.concat([].slice.call(bfv));
		}

		//console.log(txBytes.join("").toString());
			
		return txBytes.join("").toString();
		//console.log(rh)
	}

	getHash() {
		var enc = CryptoJS.SHA256(this.getBytes());
		//console.log(enc);
		//console.log(enc.toString(CryptoJS.enc.Hex));

		return enc.toString(CryptoJS.enc.Hex);
	}

	hashed() {
		this.txobj.hash = this.getHash();
	}

	getSig(priv) {
		var curve = "secp256r1";
		//var sigalg = "SHA1withECDSA";
		var sigalg = "SHA256withECDSA";
		var sig = new KJUR.crypto.Signature({"alg": sigalg});
		sig.init({d: priv, curve: curve});
		sig.updateString(this.getBytes());
		return sig.sign();
	}

	signed(priv) {
		this.txobj.sig = this.getSig(priv);
	}

	getJSON() {
		return JSON.stringify(this.txobj);
	}

	getBase64() {
		var txJSON = this.grtJSON();
		//console.log(txjson);
		return window.btoa(txJSON);
	}
	
	getSendingTx() {
		var txbase64 = this.getBase64();
		var txres = '{"call":"sendtx", "payload":"' + txbase64 + '"}';
		return txres;
	}
	
	getJSONSendingTx() {
		return window.btoa(this.getSendingTx())
	}
}