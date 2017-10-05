package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Attachment;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.attachmentRepository")
public interface AttachmentRepository extends CrudRepository<Attachment, Long> {
}
