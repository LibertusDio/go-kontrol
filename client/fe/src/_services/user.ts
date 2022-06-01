import { useLocation, useSearchParams } from "react-router-dom";
import jwt_decode from "jwt-decode";
export interface User {
    chat_id: string;
    username: string;
    first_name: string;
    last_name: string;
}
interface HomesteadToken {
    user: User;
    exp: number;
    iss: string;
    // whatever else is in the JWT.
}
export function checkUserToken() {
    // Check token query
    const loc = useLocation();
    const [searchParams] = useSearchParams();
    const { ACCESS_TOKEN_STORAGE } = process.env;
    const token = searchParams.get("token");
    if (token) {
        const user = validateToken(token);
        if (user) {
            localStorage.setItem(ACCESS_TOKEN_STORAGE ? ACCESS_TOKEN_STORAGE : "access_token", token);
            searchParams.delete("token");
            const newSearch = searchParams.toString();
            if (newSearch) {
                location.href = loc.pathname + "?" + searchParams.toString();
            } else {
                location.href = loc.pathname;
            }
            return user;
        }
        return null;
    }

    // Validate token
    const localToken = localStorage.getItem(ACCESS_TOKEN_STORAGE ? ACCESS_TOKEN_STORAGE : "access_token");
    if (localToken) {
        const user = validateToken(localToken);
        if (user) {

            return user;
        }
    }
    return null;
}

function validateToken(token: string) {
    const decodedToken: HomesteadToken = jwt_decode(token);
    if (decodedToken.exp < Date.now() / 1000) { //token expired
        return null;
    }
    return decodedToken.user;
}