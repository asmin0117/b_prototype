
var admin = require("firebase-admin");
var serviceAccount = require("./serviceAccountKey.json");
//const refreshToken = '...'; // Get refresh token from OAuth2 flow

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

module.exports = admin;