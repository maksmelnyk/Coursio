package com.example.profile.exception;

import lombok.Getter;

@Getter
public class ResourceExistsException extends RuntimeException {
    private final String errorCode;

    public ResourceExistsException(String message, String errorCode) {
        super(message);
        this.errorCode = errorCode;
    }
}