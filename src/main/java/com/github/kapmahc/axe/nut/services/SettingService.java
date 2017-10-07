package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.helper.SecurityHelper;
import com.github.kapmahc.axe.nut.models.Setting;
import com.github.kapmahc.axe.nut.repositories.SettingRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;
import java.io.*;
import java.util.Date;

@Service("nut.settingService")
@Transactional(readOnly = true)
public class SettingService {
    @Transactional(propagation = Propagation.REQUIRED)
    public void set(String k, Object v, boolean f) throws IOException {
        try (ByteArrayOutputStream bos = new ByteArrayOutputStream(); ObjectOutputStream ous = new ObjectOutputStream(bos)) {
            ous.writeObject(v);
            ous.flush();
            byte[] buf = bos.toByteArray();
            if (f) {
                buf = securityHelper.encrypt(buf);
            }
            Setting it = settingRepository.findByKey(k);
            if (it == null) {
                it = new Setting();
                it.setKey(k);
            }
            it.setValue(buf);
            it.setEncode(f);
            it.setUpdatedAt(new Date());
            settingRepository.save(it);
        }
    }

    public Object get(String k) throws IOException, ClassNotFoundException {
        Setting it = settingRepository.findByKey(k);
        if (it == null) {
            return null;
        }
        byte[] buf = it.getValue();
        if (it.isEncode()) {
            buf = securityHelper.decrypt(buf);
        }
        try (ObjectInputStream ois = new ObjectInputStream(new ByteArrayInputStream(buf))) {
            return ois.readObject();
        }
    }

    @Resource
    SettingRepository settingRepository;
    @Resource
    SecurityHelper securityHelper;
}
