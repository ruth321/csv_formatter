CREATE TABLE public.main (
  ast_nocar_trycount INTEGER,
  slyjba_taxi INTEGER,
  CONSTRAINT main_pkey PRIMARY KEY(idx),
  CONSTRAINT "main-client" FOREIGN KEY (clientid)
    REFERENCES public.client(phone)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT "main-complete" FOREIGN KEY (completeid)
    REFERENCES public.complete(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT "main-dispatchersmena" FOREIGN KEY (dispatchersmenaid)
    REFERENCES public.dispatchersmena(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT "main-rayon" FOREIGN KEY (rayonid)
    REFERENCES public.rayon(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT "main-slyjba_taxi-idx" FOREIGN KEY (slyjba_taxi)
    REFERENCES public.phonecategory(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_aautoid_fkey FOREIGN KEY (aautoid)
    REFERENCES public.auto(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_adriverid_fkey FOREIGN KEY (adriverid)
    REFERENCES public.driver(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE
    NOT VALID,
  CONSTRAINT main_avehicleid_fkey FOREIGN KEY (avehicleid)
    REFERENCES public.vehicle(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE
    NOT VALID,
  CONSTRAINT main_completedispatchersmenaid_fkey FOREIGN KEY (completedispatchersmenaid)
    REFERENCES public.dispatchersmena(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_dogovorid_fkey FOREIGN KEY (dogovorid)
    REFERENCES public.dogovor(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_offerautoid_fkey FOREIGN KEY (offerautoid)
    REFERENCES public.auto(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_orderoptionid_fkey FOREIGN KEY (orderoptionid)
    REFERENCES public.orderoption(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_orderruleid_fkey FOREIGN KEY (orderruleid)
    REFERENCES public.orderrule(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_ordersurcharge_setupid_fkey FOREIGN KEY (ordersurcharge_setupid)
    REFERENCES public.ordersurcharge_setup(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_parent_tariffid_fkey FOREIGN KEY (tariffid)
    REFERENCES public.tariff(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_phonelineid_fkey FOREIGN KEY (phonelineid)
    REFERENCES public.phoneline(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_queueid_fkey FOREIGN KEY (queueid)
    REFERENCES public.queue(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE,
  CONSTRAINT main_receivedispatchersmenaid_fkey FOREIGN KEY (receivedispatchersmenaid)
    REFERENCES public.dispatchersmena(idx)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE
) INHERITS (public.main_parent)
WITH (oids = false);

ALTER TABLE public.main
  ALTER COLUMN clientid SET STATISTICS 80;

ALTER TABLE public.main
  ALTER COLUMN createtime SET STATISTICS 80;

ALTER TABLE public.main
  ALTER COLUMN targettime SET STATISTICS 80;

ALTER TABLE public.main
  ALTER COLUMN dispatchersmenaid SET STATISTICS 80;

COMMENT ON TABLE public.main
IS 'Заказы';

COMMENT ON COLUMN public.main.idx
IS '№ заказа';

COMMENT ON COLUMN public.main.clientid
IS 'Клиент
ref client.phone/phone';

COMMENT ON COLUMN public.main.addressfrom
IS 'Откуда';

COMMENT ON COLUMN public.main.rayonid
IS 'Район';

COMMENT ON COLUMN public.main.addressto
IS 'Куда';

COMMENT ON COLUMN public.main.createtime
IS 'Создание заказа';

COMMENT ON COLUMN public.main.targettime
IS 'Время заказа';

COMMENT ON COLUMN public.main.completeid
IS 'Вид завершения';

COMMENT ON COLUMN public.main.dispatchersmenaid
IS 'Диспетчер последним изменивший заказ';

COMMENT ON COLUMN public.main.predvar
IS 'Предварительный';

COMMENT ON COLUMN public.main.comment
IS 'Комментарий';

COMMENT ON COLUMN public.main.stoimost
IS 'Cтоимость';

COMMENT ON COLUMN public.main.deleted
IS 'Удалено';

COMMENT ON COLUMN public.main.clientname
IS 'Имя клиента';

COMMENT ON COLUMN public.main.dogovor
IS 'По договору';

COMMENT ON COLUMN public.main.edittime
IS 'Время последнего редактирования заказа';

COMMENT ON COLUMN public.main.waiting
IS 'Ожидание';

COMMENT ON COLUMN public.main.name
IS 'Наименование
readonly
invisible';

COMMENT ON COLUMN public.main.queueid
IS 'Очередь автомобилей';

COMMENT ON COLUMN public.main.phonelineid
IS 'Телефонная линия';

COMMENT ON COLUMN public.main.online
IS 'Клиент на связи';

COMMENT ON COLUMN public.main.dogovorid
IS 'Договор';

COMMENT ON COLUMN public.main.receivedispatchersmenaid
IS 'Диспетчер принявший заказ';

COMMENT ON COLUMN public.main.completedispatchersmenaid
IS 'Диспетчер завершивший заказ';

COMMENT ON COLUMN public.main.aclientphone
IS 'Телефон клиента
readonly';

COMMENT ON COLUMN public.main.aautoid
IS 'Экипаж';

COMMENT ON COLUMN public.main.latitude
IS 'Широта координаты клиента';

COMMENT ON COLUMN public.main.longitude
IS 'Долгота координаты клиента';

COMMENT ON COLUMN public.main.ast_nocar_trycount
IS 'Номер попытки продлить заказ Контакт-центром';

COMMENT ON COLUMN public.main.canrejectexpl
IS 'Водитель может отказаться от заказа';

COMMENT ON COLUMN public.main.canrejectul
IS 'Спецзаказ (через адрес)';

COMMENT ON COLUMN public.main.canrejectrr
IS 'Спецзаказ (через правила распределения)';

COMMENT ON COLUMN public.main.zagorod
IS 'Загород';

COMMENT ON COLUMN public.main.addresstofull
IS 'Куда (полный адрес)';

COMMENT ON COLUMN public.main.tariffid
IS 'Тариф таксометра';

COMMENT ON COLUMN public.main.stoimost_taxometr
IS 'Стоимость поездки по таксометру';

COMMENT ON COLUMN public.main.stoimost_tarif
IS 'Стоимость поездки по тарифной сетки';

COMMENT ON COLUMN public.main.stoimost_manual
IS 'Стоимость поездки введённая вручную';

COMMENT ON COLUMN public.main.showas_freeorder
IS 'Свободный заказ';

COMMENT ON COLUMN public.main.clientid_inch
IS 'Телефон, с которого был осуществлен заказ';

COMMENT ON COLUMN public.main.orderoptionid
IS 'Услуга';

COMMENT ON COLUMN public.main.orderruleid
IS 'Правило обработки заказов';

COMMENT ON COLUMN public.main.rotateruleid
IS 'Строка цепочки распределения';

COMMENT ON COLUMN public.main.isspecnotice
IS 'Особое оповещение';

COMMENT ON COLUMN public.main.addrfrom_ulica
IS 'Откуда.Улица';

COMMENT ON COLUMN public.main.addrfrom_housenum
IS 'Откуда.Номер дома';

COMMENT ON COLUMN public.main.addrfrom_addtext
IS 'Откуда.Доп текст';

COMMENT ON COLUMN public.main.stoimost_calcroute
IS 'Стоимость поездки рассчитанная по карте';

COMMENT ON COLUMN public.main.feauteres
IS 'Характеристики требуемые от авто';

COMMENT ON COLUMN public.main.dogovor_podrid
IS 'Подразделение договора';

COMMENT ON COLUMN public.main.state
IS 'Состояние обработки заказа';

COMMENT ON COLUMN public.main.zagorod_from
IS 'Загород (откуда)';

COMMENT ON COLUMN public.main.zagorod_to
IS 'Загород (куда)';

COMMENT ON COLUMN public.main.maxcredit
IS 'Максимальная сумма безнала';

COMMENT ON COLUMN public.main.districtid
IS 'Регион';

COMMENT ON COLUMN public.main.isbonustrip
IS 'Оплата бонусами';

COMMENT ON COLUMN public.main.maxcredit_percent
IS 'Максимальный процент безнала';

COMMENT ON COLUMN public.main.rayontoid
IS 'Район КУДА';

COMMENT ON COLUMN public.main.addrlist_json
IS 'Адреса по новой версии
readonly';

COMMENT ON COLUMN public.main.latitudeto
IS 'Широта адреса Куда';

COMMENT ON COLUMN public.main.longitudeto
IS 'Долгота адреса Куда';

COMMENT ON COLUMN public.main._isnewversion
IS 'WEB-версия программы';

COMMENT ON COLUMN public.main.ordersurcharge_setupid
IS 'Поощрение для водителя';

COMMENT ON COLUMN public.main._ismanual
IS 'Редактирование диспетчером
invisible';

COMMENT ON COLUMN public.main.withcardpayment
IS 'Оплата по банк.карте';

COMMENT ON COLUMN public.main.offerautoid
IS 'Экипаж выбранный клиентом';

COMMENT ON COLUMN public.main.offerautodeadline
IS 'Конечное время для предложения водителю';

COMMENT ON COLUMN public.main.stoimost_fix
IS 'Стоимость поездки фиксированная';

COMMENT ON COLUMN public.main.orderrule_tariffid
IS 'Правило обработки заказов(строка)
ref orderrule_tariff';

COMMENT ON COLUMN public.main.paid_time
IS 'Платное время ожидания';

COMMENT ON COLUMN public.main.adrarr
IS 'Адреса (массив)';

COMMENT ON COLUMN public.main.koefsucc
IS 'Процент выполнения';

COMMENT ON COLUMN public.main.predvartime
IS 'Время предварительного заказа';

COMMENT ON COLUMN public.main.appointtime
IS 'Время назначения авто на заказ';

COMMENT ON COLUMN public.main.completetime
IS 'Выполнение заказа';

COMMENT ON COLUMN public.main.pickuptime
IS 'Плановое прибытие авто';

COMMENT ON COLUMN public.main.stoimost_features
IS 'Стоимость по характеристикам';

COMMENT ON COLUMN public.main.slyjba_taxi
IS 'Фирма/Служба такси';

COMMENT ON COLUMN public.main.couponid
IS 'Купон';

COMMENT ON COLUMN public.main.s_time_stop_taxometr
IS 'Время остановки таксометра';

COMMENT ON COLUMN public.main.bindingid
IS 'Токен карты для оплаты';

COMMENT ON COLUMN public.main.avehicleid
IS 'Автомобиль';

COMMENT ON COLUMN public.main.adriverid
IS 'Водитель';

COMMENT ON COLUMN public.main.tip
IS 'Чаевые';

CREATE INDEX "fki_main-auto" ON public.main
  USING btree (aautoid);

CREATE INDEX "fki_main-client" ON public.main
  USING btree (clientid COLLATE pg_catalog."default");

CREATE INDEX "fki_main-complete" ON public.main
  USING btree (completeid);

CREATE INDEX "fki_main-dogovor" ON public.main
  USING btree (dogovorid);

CREATE INDEX i_main_completetime ON public.main
  USING btree (completetime);

CREATE INDEX i_main_createtime ON public.main
  USING btree (createtime);

CREATE TRIGGER tg_main_oniu
  BEFORE INSERT OR UPDATE 
  ON public.main
  
FOR EACH ROW 
  EXECUTE PROCEDURE public.main_oniu();