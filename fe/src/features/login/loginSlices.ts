import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {login, logout} from "./api";
import {AuthModel, Permission} from "./authModel";
import {RootState} from "../../app/store";

const user = JSON.parse(localStorage.getItem("user") as any) as AuthModel;

export interface UserState {
    isLoggedIn: boolean
    user: Permission | string
}

const initialState: UserState = {
    isLoggedIn: !!user,
    user: user?.object_permission || ''
}
const {ACCESS_TOKEN_STORAGE} = process.env;
// @ts-ignore
const loginSlice = createSlice({
    extraReducers: {
        // @ts-ignore
        [login.fulfilled]: (state: UserState, action: PayloadAction<AuthModel>) => {
            state.isLoggedIn = action.payload.code === 200
            state.user = action.payload.object_permission
            // @ts-ignore
            localStorage.setItem("user", JSON.stringify(action.payload.object_permission))
            localStorage.setItem(ACCESS_TOKEN_STORAGE || "access_token", action.payload.object_permission.token)
        },
        // @ts-ignore
        [login.rejected]: (state: UserState, action: PayloadAction<AuthModel>) => {
            state.isLoggedIn = false
            state.user = ''
        },
        // @ts-ignore
        [login.pending]: (state: UserState, action: PayloadAction<AuthModel>) => {
            console.log("pending")
        },

        // @ts-ignore
        [logout.fulfilled]: (state: UserState, action: PayloadAction<any>) => {
            state.isLoggedIn = false
            state.user = ''
        }
    },
    initialState: initialState,
    name: "service",
    reducers: {
        clearState: (state: UserState) => {
            state.isLoggedIn = false
            state.user = ''
            localStorage.removeItem('user')
            localStorage.removeItem(ACCESS_TOKEN_STORAGE || "access_token")
            return state
        }
    }

})

export default loginSlice.reducer

export const { clearState } = loginSlice.actions;

export const userSelector = (state: RootState) => state.auth;
