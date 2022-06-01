import {createAsyncThunk} from "@reduxjs/toolkit";
import apiClient from "../../helper/apiClient";
import {AuthRequestInterface} from "./authModel";

const {ACCESS_TOKEN_STORAGE} = process.env;
export const login = createAsyncThunk("auth/login", async (payload: AuthRequestInterface, {rejectWithValue}) => {
    try {
        const resp = await apiClient().post(`/login`, payload)
        return resp.data;
    } catch (error) {
        return rejectWithValue(error);
    }
});

export const logout = createAsyncThunk("auth/logout",  () => {
    localStorage.removeItem(ACCESS_TOKEN_STORAGE || "access_token");
});
