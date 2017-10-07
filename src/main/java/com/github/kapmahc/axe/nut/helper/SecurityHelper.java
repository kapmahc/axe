package com.github.kapmahc.axe.nut.helper;

import com.google.common.primitives.Bytes;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.encrypt.BytesEncryptor;
import org.springframework.security.crypto.encrypt.Encryptors;
import org.springframework.security.crypto.encrypt.TextEncryptor;
import org.springframework.security.crypto.keygen.KeyGenerators;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import java.util.Arrays;

@Component("nut.securityHelper")
public class SecurityHelper {
    private static final int SALT_LEN = 16;

    public byte[] encrypt(byte[] plain) {
        String salt = KeyGenerators.string().generateKey();
        BytesEncryptor be = Encryptors.standard(secret, salt);
        return Bytes.concat(salt.getBytes(), be.encrypt(plain));
    }

    public byte[] decrypt(byte[] encode) {
        byte[] salt = Arrays.copyOfRange(encode, 0, SALT_LEN);
        byte[] buf = Arrays.copyOfRange(encode, SALT_LEN, encode.length);

        BytesEncryptor be = Encryptors.standard(secret, new String(salt));
        return be.decrypt(buf);
    }

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
