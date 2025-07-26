<div align="center">
  <h1>Welcome to Job Scraper!</h1>
</div>
<div align="center">
  <img src="data/job-scraper.png" alt="Image description" width="400"/>
</div>

JobScraper is a tool designed for anyone looking for jobs based on their interest, position, location (remote/on-site) using LinkedIn.
This tool can be used as both a telegram bot (access on https://t.me/jobScraperVMBot) and a CLI tool. 

Based on user input, users can either select a daily report for jobs and subscribe to email delivery or just get a report on demand as per their requirement. 

### Technical Requirements (if running on a CLI):
- MySQL installed

### Other Requirements:
- CSRF Token and Cookie from LinkedIn

### Project structure:
```
├── cmd/
│   ├── bgprocess/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── db/
├── doc/
├── data/
└── internal/
    ├── cli/
    ├── client/
    ├── core/
    │   ├── application/
    │   └── config/
    ├── models/
    ├── repository/
    ├── service/
    ├── types/
    ├── utils/
    └── workers/

```

### Setup Database (only for CLI)
Create database:
```mysql
create database job_scraper_db;
```
Create a user in mysql:
```mysql
create user 'job_scraper'@'localhost' identified by 'job_scraper';
``` 
Grant all privileges to the user
```mysql
grant all privileges on job_scraper_db.* to 'job_scraper'@'localhost';
```
Flush privileges
```mysql
flush privileges;
```

### Load the schema and pre-requisite data
```commandline
mysql -u job_scraper -p job_scraper_db < db/init.sql
mysql -u job_scraper -p job_scraper_db < db/load_prerequisites.sql
```

### Setup Postfix locally (for SMTP service)
Install postfix:
```commandline
sudo apt update
sudo apt install postfix mailutils libsasl2-modules -y
```

Update postfix config:
```commandline
sudo nano /etc/postfix/main.cf
```
and enable relay through gmail:
```text
relayhost = [smtp.gmail.com]:587
smtp_use_tls = yes
smtp_sasl_auth_enable = yes
smtp_sasl_security_options = noanonymous
smtp_sasl_password_maps = hash:/etc/postfix/sasl_passwd
smtp_tls_CAfile = /etc/ssl/certs/ca-certificates.crt
```
Create a Gmail app password (https://myaccount.google.com/security) and copy the 16 character password.

To add the gmail app password use
```commandline
sudo nano /etc/postfix/sasl_passwd
```
and append
```text
[smtp.gmail.com]:587 your-email@gmail.com:your-app-password
```
Make sure to update your-email in the above command and in the file internal/types/email.go.

Next, secure the file using
```
sudo postmap /etc/postfix/sasl_passwd
sudo chown root:root /etc/postfix/sasl_passwd /etc/postfix/sasl_passwd.db
sudo chmod 600 /etc/postfix/sasl_passwd /etc/postfix/sasl_passwd.db
```

At the end restart the postfix:
```commandline
sudo systemctl restart postfix
```

