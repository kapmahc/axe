package com.github.kapmahc.axe;

import com.github.kapmahc.axe.nut.helper.SecurityHelper;
import com.github.kapmahc.axe.nut.models.User;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import javax.annotation.Resource;

@RunWith(SpringRunner.class)
@SpringBootTest
public class ApplicationTests {

    @Test
    public void contextLoads() {
    }

    @Test
    public void testSecurity() {
        final String hello = "hello";
        System.out.printf("key=%s\n", secret);

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
    SecurityHelper securityHelper;
    @Value("${app.secret}")
    String secret;
}
