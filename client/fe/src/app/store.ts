import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import {errorHandlerMiddleware} from '../middlewares/errorMiddleware';
import errorReducer from '../_actions/errorSlices';
import loginSlice from "../features/login/loginSlices"
import serviceSlices from "../features/service/serviceSlices";

export const store = configureStore({
    reducer: {
        error: errorReducer,
        auth: loginSlice,
        service: serviceSlices
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat([
            errorHandlerMiddleware
        ]),
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType,
    RootState,
    unknown,
    Action<string>>;
