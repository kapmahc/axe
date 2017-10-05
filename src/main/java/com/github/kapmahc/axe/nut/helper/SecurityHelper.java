package com.github.kapmahc.axe.nut.helper;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.encrypt.Encryptors;
import org.springframework.security.crypto.encrypt.TextEncryptor;
import org.springframework.security.crypto.keygen.KeyGenerators;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;

@Component("nut.securityHelper")
public class SecurityHelper {
    private static final int SALT_LEN = 16;

    public String encrypt(String plain) {
        String salt = KeyGenerators.string().generateKey();
        TextEncryptor te = Encryptors.text(secret, salt);
        return salt + te.encrypt(plain);
    }

    public String decrypt(String encode) {
        TextEncryptor te = Encryptors.text(secret, encode.substring(0, SALT_LEN));
        return te.decrypt(encode.substring(SALT_LEN));
    }

    public String password(String plain) {
        return passwordEncoder.encode(plain);
    }

    public boolean check(String plain, String encode) {
        return passwordEncoder.matches(plain, encode);
    }

    @PostConstruct
    void init() {
        passwordEncoder = new BCryptPasswordEncoder();
    }

    @Value("${app.secret}")
    String secret;
    private BCryptPasswordEncoder passwordEncoder;
}
