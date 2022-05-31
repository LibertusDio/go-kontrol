import { isRejectedWithValue, Middleware } from "@reduxjs/toolkit";
import { errorAction, unauthenticationAction } from "../_actions/errorSlices";

export const errorHandlerMiddleware: Middleware = ({ dispatch }) => (next) => (action) => {
    if (isRejectedWithValue(action) && action.payload.status === 401) {
        if (action.payload.status === 401) {
            dispatch(unauthenticationAction(action.payload.data.error.message));
        } else {
            dispatch(errorAction(action.payload.data.message));
        }
    }

    return next(action);
};