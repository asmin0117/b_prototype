const mongoose = require('mongoose');
//mongodb에 connect
const Schema = mongoose.Schema;

const contractSchema = new mongoose.Schema({
  tc: {type: String, unique: true, required: true}, //계약 id
  ctype: {type: String, required: true}, // S or D
  sender: {type: Schema.Types.ObjectId, ref: 'User'}, //화물주
  receiver: {type: Schema.Types.ObjectId, ref: 'User'}, //기사
  trucksize: {type: Number},
  origin: {type: String},
  destination: {type: String},
  payment: {type: Number},
  cargosize: {type: Number}, //짐 크기
  status: {type: String}, //계약 상태
  token: {type: String}
});

module.exports = mongoose.model('Contract', contractSchema);