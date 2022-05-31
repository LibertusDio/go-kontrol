import { createSlice, PayloadAction } from "@reduxjs/toolkit";

const errorSlice = createSlice({
    name: "error",
    initialState: {
        authenticated: false,
        message: "",
    },
    reducers: {
        unauthenticationAction: (state, action: PayloadAction<string>) => {
            return {
                ...state,
                authenticated: false,
                message: action.payload,
            };
        },
        authenticationAction: (state) => {
            return {
                ...state,
                authenticated: true
            };
        },
        errorAction: (state, action: PayloadAction<string>) => {
            return {
                ...state,
                message: action.payload,
            };
        }
    }
});

export const { unauthenticationAction, authenticationAction, errorAction } = errorSlice.actions;
export default errorSlice.reducer