/*eslint-env node */
"use strict";

/**
@author: Varun Ojha
@version: 3.0
@date: 04/02/2017
@Description: SDK to talk to blockchain using hfc
**/

var Promise = require('bluebird');
var config = require('config');

var logHelper = require('../logging/logging.js');
var logger = logHelper.getLogger('blockchain_sdk');
var validate = require('../utils/validation_helper.js');
var util = require('../utils/util.js');
var constants = require('../constants/constants.js');
var bcNetwork = require('../blockchain/blockchain_network.js');


var secure = true;
var retryLimit = 5;
var retryInterval = 2000;


/**
Create a new Mortgage application
**/
function UserRegisteration(params) {
    return new Promise(function(resolve, reject){
        var UserDetails;
        try{
            logHelper.logEntryAndInput(logger, 'UserRegisteration', params);

            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'UserRegisteration', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not create mortgage application. Invalid params' })
            }

            var user = params.user;
            if(!validate.isValidString(user)){
                logHelper.logError(logger, 'UserRegisteration', 'Invalid user');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not create mortgage application. Invalid user' })
            }

            UserDetails = params.UserDetails;
          
            if(!validate.isValidJson(UserDetails)){
                logHelper.logError(logger, 'UserRegisteration', 'Invalid UserDetails');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not create mortgage application. Invalid mortgageApplication' })
            }

           
       //     var id = mortgageApplication['id'];
            var payload = JSON.stringify(UserDetails);
            logHelper.logMessage(logger,payload,'json converted to string for registeration input');
            
           
            //create an array with all values, var paylkoadArr

            var reqSpec = getRequestSpec({functionName: 'User_register', args: ["rls","12","122324244","risabh@.com","fundraiser","dsdsadad","sdfdsfs","4000000"]});
            recursiveInvoke({requestSpec: reqSpec, user: user})
            .then(function(resp){
                logHelper.logMessage(logger, 'UserRegisteration', 'Successfully registered user', resp.body);
                return resolve({statusCode: constants.SUCCESS, body: UserDetails});
            })
            .catch(function(err){   
                logHelper.logError(logger, 'UserRegisteration', 'Could not register user', err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not register user' });

            });

        }
        catch(err){
            logHelper.logError(logger, 'UserRegisteration', 'Could not register user application on blockchain ledger: ', err);
            return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not register user' });
        }
    });
}



/**
Get specific user details
**/
function getUserDetails(params) {
    console.log(params,'data in params for query method')
    return new Promise(function(resolve, reject){
       
        try{
            logHelper.logEntryAndInput(logger, 'getUserDetails', params);

            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'getUserDetails', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not fetch mortgage application. Invalid params' })
            }

            var user = params.user;
            if(!validate.isValidString(user)){
                logHelper.logError(logger, 'getUserDetails', 'Invalid user');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not fetch user details. Invalid user' })
            }

            var name = params.name;
            logHelper.logMessage(logger,name,'params.name of query method' );
            if(!validate.isValidString(name)){
                logHelper.logError(logger, 'getUserDetails', 'Invalid name');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not fetch user details. Invalid id' })
            }
                
            var reqSpec = getRequestSpec({functionName: 'readuser', args: [name]});
            recursiveQuery({requestSpec: reqSpec, user: user})
            .then(function(resp){
                logHelper.logMessage(logger, 'getUserDetails', 'Successfully fetched user details', resp.body);
                return resolve({statusCode: constants.SUCCESS, body: resp.body});
            })
            .catch(function(err){   
                logHelper.logError(logger, 'getUserDetails', 'Could not fetch user details', err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not fetch user details' });

            });

        }
        catch(err){
            logHelper.logError(logger, 'getUserDetails', 'Could not fetch property ad ', err);
            return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not fetch user details' });
        }
    });
}


/**
Generates the request object for invoke and query calls using hfc
**/
function getRequestSpec(params){

        if(!validate.isValidJson(params)){
            logHelper.logError(logger, 'getRequestSpec', 'Invalid params');
            throw new Error("Invalid params");
        }

        var chaincodeID = config['chaincode']['id'];//util.getUserDefinedProperty(constants['BLOCKCHAIN_CHAINCODE'])['id'];
        if(!validate.isValidString(chaincodeID)){
            logHelper.logError(logger, 'getRequestSpec', 'Invalid chaincodeID');
            throw new Error("Invalid chaincodeID");
        }

        var functionName = params.functionName;
        if(!validate.isValidString(functionName)){
            logHelper.logError(logger, 'getRequestSpec', 'Invalid function name');
            throw new Error("Invalid function name");
        }

        var args = []
        
        if(validate.isValidArray(params.args)){
            args = params.args;
        }

   //     var attributes = ['username', 'role']
        
     //   if(validate.isValidArray(params.attributes)){
       //     attributes = params.attributes;
        

        var spec = {
            chaincodeID: chaincodeID,
            fcn: functionName,
            args: args
          //  attrs: attributes
        }

        return spec;
}


/**
Performs query operation on blockchain
**/
function doQuery(params){
    return new Promise(function(resolve, reject){
       
        try{
            logHelper.logEntryAndInput(logger, 'doQuery', params);

            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'doQuery', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform query. Invalid params' })
            }

            var requestSpec = params.requestSpec;
            if(!validate.isValidJson(requestSpec)){
                logHelper.logError(logger, 'doQuery', 'Invalid requestSpec');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform query. Invalid requestSpec' })
            }

            var user = params.user;
            if(!validate.isValidString(user)){
                logHelper.logError(logger, 'doQuery', 'Invalid user');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform query. Invalid user' })
            }

            var chain = bcNetwork.getChain();
            var chainAsync = Promise.promisifyAll(chain);

            chainAsync.getMemberAsync(user)
            .then(function(member){
                
                var tx = member.query(requestSpec);
                tx.on('submitted', function() {
                    logHelper.logMessage(logger, 'doQuery','Transaction for query submitted');
                });

                tx.on('complete', function(data) {
                    try{
                        logHelper.logMessage(logger, 'doQuery', 'data in data.result ',data.result);
                        var buffer = new Buffer(data.result);
                         logHelper.logMessage(logger, 'doQuery', 'data in buffer.tostring ',buffer.toString());
                        var jsonResp = JSON.parse(buffer.toString());
                        logHelper.logMessage(logger,jsonResp,'json response in query line 218')
                        return resolve({statusCode: constants.SUCCESS, body: jsonResp});
                    }
                    catch(err){
                        logHelper.logError(logger,'doQuery','Could not parse query response',err);
                        return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not parse query response ' });
                    }
                });

                tx.on('error', function (err) {
                    logHelper.logError(logger, 'doQuery', 'Could not perform query ',err);
                    return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform query ' });
                   
                });
            })
            .catch(function(err){
                logHelper.logError(logger, 'doQuery', 'Could not perform query ',err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform query ' });
            })

        }
        catch(err){
                logHelper.logError(logger, 'doQuery', 'Could not perform query ',err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform query ' });
        }
    });
  
}


/**
Performs invoke operation on blockchain
**/
function doInvoke(params){
    return new Promise(function(resolve, reject){
       
        try{
            logHelper.logEntryAndInput(logger, 'doInvoke', params);

            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'doInvoke', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform invoke. Invalid params' })
            }

            var requestSpec = params.requestSpec;
            if(!validate.isValidJson(requestSpec)){
                logHelper.logError(logger, 'doInvoke', 'Invalid requestSpec');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform invoke. Invalid requestSpec' })
            }

            var user = params.user;
            if(!validate.isValidString(user)){
                logHelper.logError(logger, 'doInvoke', 'Invalid user');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not perform invoke. Invalid user' })
            }

            var chain = bcNetwork.getChain();
            var chainAsync = Promise.promisifyAll(chain);
            

            chainAsync.getMemberAsync(user)
            .then(function(member){
                
                var tx = member.invoke(requestSpec);
                tx.on('submitted', function(data) {
                    logHelper.logMessage(logger, 'doInvoke', 'Transaction for invoke submitted ',requestSpec);
                    return resolve({statusCode: constants.SUCCESS, body: data});
                    
                });

                tx.on('complete', function(data) {
                    //logHelper.logMessage(logger, 'doInvoke', 'Transaction for invoke complete ',data);
                    
                });

                tx.on('error', function (err) {
                    logHelper.logError(logger, 'doInvoke', 'Could not perform invoke ',err);
                    return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform invoke ' });
                   
                });
            })
            .catch(function(err){
                logHelper.logError(logger, 'doInvoke', 'Could not perform invoke ',err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform invoke ' });
            })

        }
        catch(err){
                logHelper.logError(logger, 'doInvoke', 'Could not perform invoke ',err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not perform invoke ' });
        }
    });
  
}

/*Performs register on Blockchain CA*/
function doRegister(params) {
    return new Promise(function(resolve, reject){
        
        var username;
        try{
            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'registerUser', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not register user. Invalid params' })
            }

            username = params.username;
            if(!validate.isValidString(username)){
                logHelper.logError(logger, 'registerUser', 'Invalid username');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not register user. Invalid username' })
            }

            var affiliation = params.affiliation;
            if(!validate.isValidString(affiliation)){
                logHelper.logError(logger, 'registerUser', 'Invalid affiliation');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not register user. Invalid affiliation' })
            }

            var roles = params.roles;
            if(!validate.isValidArray(roles)){
                roles = ['client'];
            }

            var enrollsecret
            var chain = bcNetwork.getChain();
            var reg = chain.getRegistrar();
            var chainAsync = Promise.promisifyAll(chain);

            chainAsync.getMemberAsync(username)
            .then(function(member){
                var memberAsync = Promise.promisifyAll(member);
                
                    var registrationRequest = {
                        enrollmentID: username,
                        attributes: [
                            {name: 'role', value: affiliation},
                            {name: 'username', value: username}
                        ],
                        affiliation: 'group1',
                        registrar: reg,
                        roles: roles
                        
                    };
                    
                return memberAsync.registerAsync(registrationRequest);
            })
            .then(function(enrollsec){
                logHelper.logMessage(logger, 'registerUser', 'Successfully registered user on blockchain: '+username);
                enrollsecret = enrollsec;
                return resolve({statusCode: constants.SUCCESS, body: {password: enrollsecret}});
                
            })
            .catch(function(err){
                logHelper.logError(logger, 'registerUser', 'Could not register user on blockchain: '+username, err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not register user' });
            })
        }
        catch(err){
            logHelper.logError(logger, 'registerUser', 'Could not register user on blockchain: '+username, err);
            return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not register user' });
        }
    });
   
}

/**
Enroll user with the Blockchain CA
**/
function doLogin(params) {
    return new Promise(function(resolve, reject){

        try{
            logHelper.logMethodEntry(logger, 'doLogin');
            if(!validate.isValidJson(params)){
                logHelper.logError(logger, 'doLogin', 'Invalid params');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not login user. Invalid params' })
            }

            var username = params.username;
            if(!validate.isValidString(username)){
                logHelper.logError(logger, 'doLogin', 'Invalid username');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not login user. Invalid username' })
            }

            var password = params.password;
            if(!validate.isValidString(password)){
                logHelper.logError(logger, 'doLogin', 'Invalid account');
                return reject({statusCode: constants.INVALID_INPUT, body: 'Could not login user. Invalid password' })
            }

            var chain = bcNetwork.getChain();
            var chainAsync = Promise.promisifyAll(chain);

            chainAsync.getMemberAsync(username)
            .then(function(member){
                var memberAsync = Promise.promisifyAll(member);
                return memberAsync.enrollAsync(password);
            })
            .then(function(crypto){
                logHelper.logMessage(logger, 'doLogin', 'Successfully logged in user on blockchain: '+username);
                return resolve({statusCode: constants.SUCCESS, body: crypto});
            })
            .catch(function(err){
                logHelper.logError(logger, 'doLogin', 'Could not login user on blockchain: '+username, err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not login user' });
            });

        }
        catch(err){
            logHelper.logError(logger, 'doLogin', 'Could not register user on blockchain: '+username, err);
            return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not register user' });
        }
    });
}



function recursiveRegister(params){
    if(!validate.isValid(params.retryCounter)){
        params.retryCounter = 0;
    }
    else{
        params.retryCounter = params.retryCounter + 1;
    }

    return doRegister(params).catch(function(err) {
        if(err.statusCode == constants.INVALID_INPUT || params.retryCounter > retryLimit){
            logHelper.logError(logger, 'recursiveRegister', "Register Retries Exhausted", err)
            return Promise.reject(err);
        }
        return Promise.delay(retryInterval).then(function(){
            return recursiveRegister(params);
        });
    });
}

function recursiveLogin(params){
    if(!validate.isValid(params.retryCounter)){
        params.retryCounter = 0;
    }
    else{
        params.retryCounter = params.retryCounter + 1;
    }

    return doLogin(params).catch(function(err) {
        if(err.statusCode == constants.INVALID_INPUT || params.retryCounter > retryLimit){
            logHelper.logError(logger, 'recursiveLogin', "Login Retries Exhausted", err)
            return Promise.reject(err);
        }
        return Promise.delay(retryInterval).then(function(){
            return recursiveLogin(params);
        });
    });
}



function recursiveInvoke(params){
    if(!validate.isValid(params.retryCounter)){
        params.retryCounter = 0;
    }
    else{
        params.retryCounter = params.retryCounter + 1;
    }

    return doInvoke(params).catch(function(err) {
        if(err.statusCode == constants.INVALID_INPUT || params.retryCounter > retryLimit){
            logHelper.logError(logger, 'recursiveInvoke', "Invoke Retries Exhausted", err);
            return Promise.reject(err);
        }
        return Promise.delay(retryInterval).then(function(){
            logHelper.logError(logger, 'recursiveInvoke', "Invoke Retry "+params.retryCounter, err)
            return recursiveInvoke(params);
        });
    });
}

function recursiveQuery(params){
   if(!validate.isValid(params.retryCounter)){
        params.retryCounter = 0;
    }
    else{
        params.retryCounter = params.retryCounter + 1;
    }

    return doQuery(params).catch(function(err) {
        if(err.statusCode == constants.INVALID_INPUT || params.retryCounter > retryLimit){
            logHelper.logError(logger, 'recursiveQuery', "Query Retries Exhausted", err)
            return Promise.reject(err);
        }
        return Promise.delay(retryInterval).then(function(){
            logHelper.logError(logger, 'recursiveQuery', "Query Retry "+params.retryCounter, err)
            return recursiveQuery(params);
        });
    });
}

function isUserRegistered(params){
    return new Promise(function(resolve,reject){
        try{
            logHelper.logMethodEntry(logger,'isUserRegistered');
            var username = params.username;
            var chain = bcNetwork.getChain();
            var chainAsync = Promise.promisifyAll(chain);

            chainAsync.getMemberAsync(username)
            .then(function(member){
                return resolve({statusCode: constants.SUCCESS, body: member.isRegistered()});
            })
            .catch(function(err){
                logHelper.logError(logger, 'isUserRegistered', 'Could not get user registration status '+username, err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not get user registration status' });
            });
        }
        catch(err){
             logHelper.logError(logger, 'isUserRegistered', 'Could not get user registration status '+username, err);
             return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not get user registration status' });
        }
    })
   
}

function isUserEnrolled(params){
    return new Promise(function(resolve,reject){
        try{
            logHelper.logMethodEntry(logger,'isUserEnrolled');
            var username = params.username;
            var chain = bcNetwork.getChain();
            var chainAsync = Promise.promisifyAll(chain);

            chainAsync.getMemberAsync(username)
            .then(function(member){
                return resolve({statusCode: constants.SUCCESS, body: member.isEnrolled()});
            })
            .catch(function(err){
                logHelper.logError(logger, 'isUserEnrolled', 'Could not get user enrollment status '+username, err);
                return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not get user enrollment status' });
            });
        }
        catch(err){
             logHelper.logError(logger, 'isUserEnrolled', 'Could not get user enrollment status '+username, err);
             return reject({statusCode: constants.INTERNAL_SERVER_ERROR, body: 'Could not get user enrollment status' });
        }
    })
   
}

module.exports = {
    UserRegisteration: UserRegisteration,
    getUserDetails: getUserDetails,
    recursiveRegister: recursiveRegister,
    recursiveLogin: recursiveLogin,
    isUserEnrolled: isUserEnrolled,
    isUserRegistered: isUserRegistered
}

