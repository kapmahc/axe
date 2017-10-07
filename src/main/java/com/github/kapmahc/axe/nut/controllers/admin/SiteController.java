package com.github.kapmahc.axe.nut.controllers.admin;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import javax.annotation.Resource;
import javax.sql.DataSource;
import java.net.*;
import java.sql.DatabaseMetaData;
import java.sql.SQLException;
import java.util.*;

@Controller("auth.admin.siteController")
@RequestMapping(value = "/admin/site")
public class SiteController {
    @GetMapping("/status")
    public String getStatus(Model model) throws SQLException, SocketException, UnknownHostException {
        model.addAttribute("os", osStatus());
        model.addAttribute("java", javaStatus());
        model.addAttribute("network", networkStatus());
        model.addAttribute("redis", redisStatus());
        model.addAttribute("database", databaseStatus());
        model.addAttribute("rabbitmq", rabbitmqStatus());
        return "nut/admin/site/status";
    }

    private Map<Object, Object> redisStatus() {
        Map<Object, Object> map = new LinkedHashMap<>();
        redisTemplate.executeWithStickyConnection((c) -> {
            Properties props = c.info();
            props.forEach(map::put);
            return null;
        });
        return map;
    }

    private Map<String, Object> rabbitmqStatus() {
        Map<String, Object> map = new LinkedHashMap<>();
        rabbitTemplate.execute((ch) -> map.put("provider", ch.getConnection().getClientProvidedName()));
        return map;
    }

    private Map<String, Object> databaseStatus() throws SQLException {
        Map<String, Object> map = new LinkedHashMap<>();
        DatabaseMetaData meta = dataSource.getConnection().getMetaData();
        map.put(
                "product",
                String.format(
                        "%s %s",
                        meta.getDatabaseProductName(),
                        meta.getDatabaseProductVersion()
                )
        );
        map.put(
                "driver",
                String.format(
                        "%s %s",
                        meta.getDriverName(),
                        meta.getDriverVersion()
                )
        );
        map.put("url", String.format("%s@%s", meta.getUserName(), meta.getURL()));
        return map;
    }

    private Map<String, Object> networkStatus() throws UnknownHostException, SocketException {
        Map<String, Object> map = new LinkedHashMap<>();
        InetAddress ia = InetAddress.getLocalHost();
        map.put("name", ia.getHostName());
        map.put("address", ia.getHostAddress());

        List<String> eth = new ArrayList<>();
        Enumeration<NetworkInterface> b = NetworkInterface.getNetworkInterfaces();
        while (b.hasMoreElements()) {
            for (InterfaceAddress f : b.nextElement().getInterfaceAddresses())
                if (f.getAddress().isSiteLocalAddress()) {
                    eth.add(f.getAddress().getHostAddress());
                }
        }
        map.put("interfaces", eth);
        return map;
    }

    private Map<String, Object> javaStatus() {

        Map<String, Object> map = new LinkedHashMap<>();
        for (String k : new String[]{
                "home", "version", "vendor", "vendor.url"
        }) {
            map.put(k, System.getProperty("java." + k));
        }
        Runtime run = Runtime.getRuntime();
        map.put(
                "memory usage",
                String.format(
                        "free(%dMB) total(%dMB) max(%dMB)",
                        run.freeMemory() / 1024 / 1024,
                        run.totalMemory() / 1024 / 1024,
                        run.maxMemory() / 1024 / 1024
                )
        );
        Package pkg = getClass().getPackage();
        map.put(
                "application",
                String.format(
                        "%s (%s)",
                        pkg.getImplementationTitle(),
                        pkg.getImplementationVersion()
                )
        );
        map.put("classpath", System.getProperty("java.class.path").replaceAll(":", "\n"));
        return map;
    }

    private Map<String, Object> osStatus() {
        Map<String, Object> map = new LinkedHashMap<>();
        for (String k : new String[]{
                "name", "arch", "version"
        }) {
            map.put(k, System.getProperty("os." + k));
        }
        map.put(
                "user",
                String.format(
                        "%s@%s",
                        System.getProperty("user.name"),
                        System.getProperty("user.dir")
                )
        );


        map.put("cpu cores", Runtime.getRuntime().availableProcessors());

        return map;
    }

    @Resource
    DataSource dataSource;
    @Resource
    RedisTemplate redisTemplate;
    @Resource
    RabbitTemplate rabbitTemplate;
//    @Resource
//    RabbitManagementTemplate rabbitManagementTemplate;


}
