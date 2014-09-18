## Purpose

This project is a [Go](http://www.golang.org) implementation of the [MangoPay HTTP REST api](http://www.mangopay.com/) version 2.

## Installation

Use the api with

    $ go get github.com/gotsunami/mangopay2-go-sdk

A CLI tool is also available for testing the MangoPay service easily:

    $ go get github.com/gotsunami/mangopay2-go-sdk/mangopay

Now run `mangopay` from a terminal to get the list of supported actions:

    Usage: ./mangopay [options] action configfile
     
    where action is one of: 
      conf              show config
      events            list all events (PayIns, PayOuts, Transfers)
      users             list all users
      user*             fetch a user (natural or legal)

      addnatuser*       create a natural user
      editnatuser*      update natural user info
      natuser*          fetch natural user info

      addlegaluser*     create a legal user
      editlegaluser*    update legal user info
      legaluser*        fetch legal user info

      addwallet*        create a new wallet
      editwallet*       update wallet info
      trwallet*         fetch all wallet's transactions
      wallet*           fetch wallet info
      wallets*          list all user's wallet

      addtransfer*      create a new tranfer
      transfer*         fetch transfer info
      transfers*        list all user's transactions

      addwebpayin*      create a payIn through web interface
      adddirectpayin*   create a direct payIn (with tokenized card)
      payin*            fetch a payIn

      addcard*          register a credit card
      card*             fetch a credit card
      cards*            list all user's cards

      addrefund*        refund a payment (provide TransferId or PayInId)
      refund*           fetch a refund (transfer or payin)

      addaccount*       create an IBAN bank account
      account*          fetch a user's bank account
      accounts*         list all user's bank accounts

      addpayout*        create a bank wire
      payout*           fetch a bank wire

    Actions with an asterisk(*) require input JSON data (-d).

    Options:
      -d="": JSON data part of the HTTP request
      -v=0: Verbosity level (1 for debug)
