import {getKindeServerSession} from "@kinde-oss/kinde-auth-nextjs/server";
import { prismadb } from "@/lib/prismadb";
import { NextResponse } from "next/server";

/* The purpose of this endpoint is to do the following:

1. Given a kinde usersession, get the current user
2. If the user does not exist in Kinde, the we return an unauthorized response
3. If the user exists in Kinde, we check if it exists in the database
4. If the user does not exist in the database, we create it
5. If the user exists in the database, we return a success response

*/

export async function GET(params: Request) {

    // get user from the current session
    const { getUser } = getKindeServerSession();
    const user = await getUser();

    // User does not exist return unauthorized
    if (!user || !user.id ) return NextResponse.json({ message: "Unauthorized", success: false}, {status: 401});

    
    try{
        // check if user in  the database
        const dbUser = await prismadb.user.findUnique({
            where: {id: user.id}
        });

        // if not create user
        if (!dbUser) {
        
        const newdbUser = await prismadb.user.create({
                data: {
                    userName: `${user.given_name}_${user.id.toString().substring(0, 4)}`,
                    email: user.email,
                }
            })

            return NextResponse.json({ message:`User Id ${newdbUser.id} synced to database`, success: true, user: newdbUser}, {status: 200});
        }
        return NextResponse.json({ message: "User already exists", success: true, user: dbUser}, {status: 200});
       

    } catch (error) {
        console.log(error)
    }
    
}