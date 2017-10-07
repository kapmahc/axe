package com.github.kapmahc.axe;

import com.github.kapmahc.axe.nut.helper.JwtHelper;
import com.github.kapmahc.axe.nut.helper.SecurityHelper;
import com.github.kapmahc.axe.nut.models.Log;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.services.SettingService;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import javax.annotation.Resource;
import java.io.IOException;
import java.time.Duration;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

@RunWith(SpringRunner.class)
@SpringBootTest
public class ApplicationTests {


    @Test
    public void testSetting() {
        Log it = new Log();
        it.setId(123456L);
        it.setMessage(hello);
        it.setCreatedAt(new Date());

        final String k1 = "test.k1";
        final String k2 = "test.k2";
        try {
            settingService.set(k1, it, false);
            settingService.set(k2, it, true);

            Log v1 = (Log) settingService.get(k1);
            Log v2 = (Log) settingService.get(k2);
            System.out.printf("%s\n%s\n%s\n", it.getCreatedAt(), v1.getCreatedAt(), v2.getCreatedAt());
            assert it.getId().equals(v1.getId()) && it.getId().equals(v2.getId());
            assert hello.equals(v1.getMessage()) && hello.equals(v2.getMessage());
        } catch (IOException | ClassNotFoundException e) {
            throw new RuntimeException(e);
        }

    }


    @Test
    public void testJwt() {
        final String key = "hi";
        Map<String, String> claim = new HashMap<>();
        claim.put(key, hello);
        String token = jwtHelper.generate(claim, Duration.ofMinutes(1));
        System.out.printf("jwt token = %s\n", token);
        assert hello.equals(jwtHelper.parse(token).get(key));
    }

    @Test
    public void testSecurity() {
        String passwd = securityHelper.password(hello);
        System.out.printf("password(%s) = %s\n", hello, passwd);
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
    @Resource
    SettingService settingService;
    @Value("${app.secret}")
    String secret;
    private final String hello = "Hello, AXE!";
}
