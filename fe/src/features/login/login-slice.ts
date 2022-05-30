import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import LoginApi from "./login-api";
import {AuthModel, AuthRequestInterface, Permission} from "./AuthModel";
import {RootState} from "../../app/store";

const user = JSON.parse(localStorage.getItem("user") as any) as AuthModel;

export const login = createAsyncThunk(
    "auth/login",
    async (u: AuthRequestInterface, thunkAPI: any) => {
        try {
            const data = await LoginApi.login(u);
            return {user: data};
        } catch (error) {
            // thunkAPI.dispatch("xxxxxxx");
            return thunkAPI.rejectWithValue("xxxxxx");
        }
    }
);
export const logout = createAsyncThunk("auth/logout", async () => {
    await LoginApi.logout();
});

export interface UserState {
    isLoggedIn: boolean
    user: any
}

const initialState: UserState = {
    isLoggedIn: !!user,
    user: user || ''
}

// @ts-ignore
const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        clearState: (state: UserState) => {
            state.isLoggedIn = false
            state.user = ''
            localStorage.removeItem("token")
            localStorage.removeItem("user")
            return state
        }
    },
    extraReducers: {
        // @ts-ignore
        [login.fulfilled]: (state: UserState, { p} : AuthModel) => {
            state.isLoggedIn = true
            state.user = p as AuthModel
            // state.isFetching = false;
            // state.isSuccess = true;
            // state.email = payload.user.email;
            // state.username = payload.user.name;
        },
        // [register.fulfilled]: (state, action) => {
        //     state.isLoggedIn = false;
        // },
        // [register.rejected]: (state, action) => {
        //     state.isLoggedIn = false;
        // },
        // [login.fulfilled]: (state, action) => {
        //     state.isLoggedIn = true;
        //     state.user = action.payload.user;
        // },
        // [login.rejected]: (state, action) => {
        //     state.isLoggedIn = false;
        //     state.user = null;
        // },
        // [logout.fulfilled]: (state, action) => {
        //     state.isLoggedIn = false;
        //     state.user = null;
        // },
    },
});
export default authSlice.reducer;
export const { clearState } = authSlice.actions;
export const userSelector = (state: RootState) => state.auth;
