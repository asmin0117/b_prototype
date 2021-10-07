// server.js
// where your node app starts

const express = require("express");
const jwt = require("jsonwebtoken");
const bodyParser = require("body-parser");
const kakao_auth = require("./kakao_auth.js");
require("dotenv").config();

var mongoose = require("mongoose");
var db = mongoose.connection;
db.on("error", console.error);
db.once("open", function () {
  console.log("connected to db");
});
mongoose.connect(process.env.mongoURI);

//User
const User = require("./model/user");

const bycrypt = require("bcryptjs");

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "connection.json");
const ccpJSON = fs.readFileSync(ccpPath, "utf8");
const ccp = JSON.parse(ccpJSON);

const app = express();

//app.use(bodyParser.urlencoded({ extended: false }));
app.use(express.json({ extended: false }));

// make all the files in 'public' available
// https://expressjs.com/en/starter/static-files.html
app.use(express.static("public"));

// https://expressjs.com/en/starter/basic-routing.html
app.get("/", (request, response) => {
  response.sendFile(__dirname + "/views/index.html");
});

app.get("/callbacks/kakao/login", async (request, response) => {
  //Auth code 받아서 돌려주는 API
  try {
    const redirect = `webauthcallback://success?${new URLSearchParams(request.query).toString()}`;
    console.log(`Redirecting to ${redirect}`);
    response.redirect(307, redirect);
  } catch (error) {
    print(error);
  }
  
});

app.post("/callbacks/kakao/token", async (request, response) => {
  //발급 받은 kakao AccessCode로 사용자 확인 후 firebase로 custom token 생성 API
  //console.log(request.body["accessToken"]);
  try {
    await kakao_auth.createFirebaseToken(request.body["accessToken"], (result) => {
      esponse.send(result);
    });
  } catch (error) {
    print(error);
  }
});

// 라우팅
app.post(`/Signup`, async (req, res) => {
  const { email, name, phone, account } = req.body;
  console.log(email, name, phone, account);

  // db작업
  // db결과확인 -> token id, name
  try {
    let user = await User.findOne({ phone });
    console.log(user);
    if (user) {
      return res.status(400).json({ errors: [{ msg: "이미 존재하는 사용자입니다." }] });
    }

    user = new User({
      email,
      name,
      phone,
      account,
    });

    //const salt = await bycrypt.genSalt(10);
    //user.account = await bycrypt.hash(account, salt);

    await user.save();

    // cc -> regPERSON 호출 ( name, token )
    result = await cc_call("regPERSON", [name, account, "0"]);
    const s_message = { result: "cc 호출, regPERSON" };
    res.status(200).json({ msg: s_message, customMsg: "사용자 등록 완료" });
  } catch (error) {
    console.error(error);
    res.status(500).send("500 에러");
  }
});

async function cc_call(fn_name, args) {
  const walletPath = path.join(process.cwd(), "wallet");
  const wallet = new FileSystemWallet(walletPath);

  const userExists = await wallet.exists("user1");
  if (!userExists) {
    console.log('An identity for the user "user1" does not exist in the wallet');
    console.log("Run the registerUser.js application before retrying");
    return;
  }
  const gateway = new Gateway();
  await gateway.connect(ccp, { wallet, identity: "user1", discovery: { enabled: false, asLocalhost: true } });
  const network = await gateway.getNetwork("mychannel");
  const contract = network.getContract("cargo");

  var result;

  if (fn_name == "regPERSON") {
    result = await contract.submitTransaction("regPERSON", args[0], args[1], args[2]);
  } else if (fn_name == "regTransportation") {
    result = await contract.submitTransaction("regTransportation", args[0], args[1], args[2], args[3]);
  } else if (fn_name == "history") {
    result = await contract.evaluateTransaction("history", args[0]);
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
  } else {
    result = "not supported function";
  }

  gateway.disconnect();

  return result;
}

// listen for requests :)
const listener = app.listen(process.env.port, () => {
  console.log("Your app is listening on port " + listener.address().port);
});
