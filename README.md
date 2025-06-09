# Job Scraper
#### Welcome to Job Scraper CLI tool. This tool is for everyone who wants to find jobs based on their interest, position, location (remote/on-site) using online job platforms like LinkedIn, Glassdoor, Indeed, etc using their authtokens.

### Requirements:
- Authtokens from Linkedin, glassdoor, indeed, etc

### Project structure:
```
cmd/
├── api/
│   └── main.go
|
internal/
├── core/
│   ├── middlewares/
│   ├── config/     
│   ├── database/   
│   ├── application/
├── common/
│   ├── utils/
│   ├── log/
|   ├── types/
├── routes/
├── controllers/
├── models/
├── repositories/
├── services/
└── workers/
```