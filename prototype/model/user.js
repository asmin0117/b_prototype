const mongoose = require("mongoose");

const userSchema = new mongoose.Schema({
  email: { type: String, unique: true },
  name: { type: String },
  account: { type: String, unique: true }, //UserCharID
  phone: { type: String, unique: true, required: true },
  token: { type: String }
});

module.exports = mongoose.model("user", userSchema);
