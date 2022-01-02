Building a Simple Pay-Later Service

As a pay later service we allow our users to buy goods from a merchant now, and then allow
them to pay for those goods at a later date.
The service works inside the boundary of following simple constraints -
1. Let's say that for every transaction paid through us, merchants offer us a discount.
○ For example, if the transaction amount is Rs.100, and merchant discount offered to us is 10%, we pay Rs. 90 back to the merchant.
2. The discount varies from merchant to merchant.
3. A merchant can decide to change the discount it offers to us, at any point in time.
4. All users get onboarded with a credit limit, beyond which they can't transact.
5. If a transaction value crosses this credit limit, we reject the transaction.

Use Cases
There are various use cases our service is intended to fulfil -
● allow merchants to be onboarded with the amount of discounts they offer
● allow merchants to change the discount they offer
● allow users to be onboarded (name, email-id and credit-limit)
● allow a user to carry out a transaction of some amount with a merchant.
● allow a user to pay back their dues (full or partial)

Reporting:
● how much discount we received from a merchant till date
● dues for a user so far
● which users have reached their credit limit
● total dues from all users together

Goal
The goal of this coding challenge will be to build a system for satisfying above use cases.
● IO will be via a command line interface.
● The input can be given in any order as a command, and the system should respond
accordingly.
● For inputs like merchant discount rate changes or credit limit changes for a user, the
system adapts itself.

CLI
here is how the command line interface, corresponding to the use-cases mentioned above, can look like -

Use the command "simpl" in the first

new user u1 u1@email.in 1000 # name, email, credit-limit
new merchant m1 2%           # name, discount-percentage
new txn u1 m2 400            # user, merchant, txn-amount
update merchant m1 1%        # merchant, new-discount-rate
payback u1 300               # user, payback-amount
report discount m1
report dues u1
report users-at-credit-limit
report total-dues

To create the database/tables:

mysql query
1. CREATE DATABASE simpl;
2. CREATE TABLE users(  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    user_name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    credit_limit int,
    spent int
);

3. CREATE TABLE merchants (  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    merchant_name VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    discount FLOAT,
    total_discount FLOAT
);

Basic commands 
1. simpl -h
2. simpl new -h
3. simpl update -h
4. simpl report -h
5. simpl payback -h


Example Flow
> simpl new user user1 u1@users.com 300
user1(300)
> simpl  new user user2 u2@users.com 400
user2(400)
> simpl  new user user3 u3@users.com 500
user3(500)
> simpl  new merchant m1 m1@merchants.com 0.5%
m1(0.5%)
> simpl  new merchant m2 m2@merchants.com 1.5%
m2(1.5%)
> simpl  new merchant m3 m3@merchants.com 1.25%
m3(1.25%)
> simpl  new txn user2 m1 500
rejected! (reason: credit limit)
> simpl  new txn user1 m2 300
success!
> simpl  new txn user1 m3 10
rejected! (reason: credit limit)
> simpl  report users-at-credit-limit
user1
> simpl  new txn user3 m3 200
success!
> simpl  new txn user3 m3 300
success!
> simpl  report users-at-credit-limit
user1
user3
> simpl  report discount m3
6.25
> simpl  payback user3 400
user3(dues: 100)
> simpl  report total-dues
user1: 300
user3: 100
total: 400