package com.github.kapmahc.axe.ops.vpn.repositories;

import com.github.kapmahc.axe.ops.vpn.models.Log;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.vpn.LogRepository")
public interface LogRepository  extends CrudRepository<Log, Long> {
}
