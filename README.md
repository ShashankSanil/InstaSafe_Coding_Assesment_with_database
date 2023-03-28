# InstaSafe_Coding_Assesment_with_databse
Connected with MongoDB database.

# Run the code by using the following commands :

 --> cd server
 
 
 --> go run main.go

# API

--> Import "PostmanCollection.json" file to Postman.

1-> Create_End_User : POST
    
    url: http://localhost:5055/user/signUp

    sampleInput :   {
                        "Username":"Shashank",
                        "Email":"Shashank@gmail.com"
                    }


2-> Add_Transaction : POST

    url: http://localhost:5055/transactions

    sampleInput :   {
                        "amount":"10.23",
                        "end_user_name":"Shashank",
                        "end_user_email":"Shashank@gmail.com",
                        "timestamp":"2023-03-27T16:14:22.194962Z",
                        "location":"Bangalore"
                    }


3-> Add_Loaction : POST

    url: http://localhost:5055/user/:uid/addLoaction  

         (http://localhost:5055/user/6421bfb857da685452441912/addLoaction)

    sampleInput :   {
                        "city":"Mangalore"
                    }

4-> Reset_Loaction : PUT

    url: http://localhost:5055/user/:uid/resetLoaction

        (http://localhost:5055/user/6421bfb857da685452441912/resetLoaction)

    sampleInput :   {
                        "city":"Udupi"
                    }

5-> Get_Statistics : GET

    url: http://localhost:5055/user/:uid/statistics

         (http://localhost:5055/user/6421bfb857da685452441912/statistics?city=Mangalore)


6-> Delete_All_Transactions : DELETE

    url: http://localhost:5055/transactions

  