import axios from "axios";
import {AuthModel, AuthRequestInterface} from "./AuthModel";
import API_URL from '../../constants';

const login = (payload: AuthRequestInterface): AuthModel | any => {
    return axios
        .post(API_URL + "/login", payload)
        .then((response: any) => {
            if (response.data) {
                localStorage.setItem("user", JSON.stringify(response.data));
            }
            const authModel = response.data as AuthModel
            localStorage.setItem('token', authModel.object_permission.token)
            return authModel;
        });
};
const logout = () => {
    localStorage.removeItem("user");
};
const authService = {
    login,
    logout,
};
export default authService;
