package com.github.kapmahc.axe.nut.helper;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTCreator;
import com.auth0.jwt.JWTVerifier;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import java.io.UnsupportedEncodingException;
import java.time.Duration;
import java.time.ZoneId;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

@Component("auth.jwtHelper")
public class JwtHelper {
    public Map<String, String> parse(String token) {
        JWTVerifier verifier = JWT.require(algorithm)
                .withIssuer(name)
                .build();
        DecodedJWT jwt = verifier.verify(token);
        Map<String, String> map = new HashMap<>();
        jwt.getClaims().forEach((k, v) -> map.put(k, v.asString()));
        return map;
    }

    public String generate(Map<String, String> claims, Duration dur) {
        Date now = new Date();
        ZoneId zone = ZoneId.systemDefault();
        Date end = Date.from(
                now.toInstant()
                        .atZone(zone)
                        .toLocalDateTime()
                        .plusSeconds(dur.getSeconds())
                        .atZone(zone)
                        .toInstant()
        );
        JWTCreator.Builder builder = JWT.create()
                .withIssuer(name)
                .withNotBefore(now)
                .withExpiresAt(end);
        claims.forEach(builder::withClaim);
        return builder.sign(algorithm);
    }

    @PostConstruct
    void init() throws UnsupportedEncodingException {
        algorithm = Algorithm.HMAC512(secret);
    }

    @Value("${app.secret}")
    String secret;
    @Value("${app.name}")
    String name;
    private Algorithm algorithm;
}
