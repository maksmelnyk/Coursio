package com.example.profile.exceptions;

public class AuthenticationException extends BaseAppException {
    public AuthenticationException(String message, String errorCode) {
        super(message, errorCode);
    }
}