package com.github.kapmahc.axe;

import com.github.kapmahc.axe.nut.helper.JwtHelper;
import com.github.kapmahc.axe.nut.helper.SecurityHelper;
import com.github.kapmahc.axe.nut.models.User;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import javax.annotation.Resource;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

@RunWith(SpringRunner.class)
@SpringBootTest
public class ApplicationTests {


    @Test
    public void contextLoads() {
    }

    @Test
    public void testJwt() {
        final String key = "hi";
        Map<String, String> claim = new HashMap<>();
        claim.put(key, hello);
        String token = jwtHelper.generate(claim, Duration.ofMinutes(1));
        System.out.printf("jwt token=%s\n", token);
        assert hello.equals(jwtHelper.parse(token).get(key));
    }

    @Test
    public void testSecurity() {
        String passwd = securityHelper.password(hello);
        System.out.printf("password(%s)=%s\n", hello, passwd);
        assert securityHelper.check(hello, passwd);

        String encode = securityHelper.encrypt(hello);
        System.out.printf("encrypt(%s) = %s\n", hello, encode);
        assert hello.equals(securityHelper.decrypt(encode));
    }

    @Test
    public void testGravatar() throws Exception {
        User it = new User();
        it.setEmail(" MyEmailAddress@example.com ");
        it.setGravatarLogo();
        assert "0bc83cb571cd1c50ba6f3e8a78ef1346".equals(it.getLogo());
    }

    @Resource
    JwtHelper jwtHelper;
    @Resource
    SecurityHelper securityHelper;
    @Value("${app.secret}")
    String secret;
    private final String hello = "Hello, AXE!";
}
