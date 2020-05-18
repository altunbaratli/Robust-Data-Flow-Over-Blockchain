#!/usr/bin/env python
# coding: utf-8

# In[2]:


from pyspark import SparkContext, SQLContext, SparkConf
from pyspark.sql import SparkSession
from pyspark.sql import Row
from pyspark.sql import functions as F
from pyspark.sql.functions import from_utc_timestamp
import os

# spark.stop()

conf = SparkConf().setAppName("PySpark App").setMaster("local").set("spark.jars", "C:\\Program Files\\PostgreSQL\\12\\postgresql-42.2.8.jar").set('spark.driver.host','127.0.0.1')
sc = SparkContext(conf=conf)
sqlContext = SQLContext(sc)
spark = SparkSession     .builder     .appName("Python Spark SQL basic example").enableHiveSupport().getOrCreate()


# In[7]:


import requests
pload = {'username':'user3','channel':'channel1', 'smartcontract':'cc', 'args': {'sensorID':'1'}}
r = requests.post('http://173.193.75.191:30006/api/getHistory',json = pload)
print(r.text)


# In[8]:


json_rdd = sc.parallelize([r.text])
convertdf = spark.read.json(json_rdd)
convertdf.printSchema()


# In[13]:


from pyspark.sql import functions as F

newdf = convertdf.select("Value.sensorID", "Value.time", "Value.temp", "Value.stuckatfault", "Value.scfault", "Value.outlier")
flattened = newdf.dropna()
flattened_full = flattened.where("temp!=' '")
# df = flattened_full.withColumn('new_date',F.to_date(F.unix_timestamp('time', 'yyyy-MM-dd HH:mm:ss').cast('timestamp')))
df = flattened_full.withColumn('temp',F.col('temp').cast('float'))
df.toPandas().plot(x="time", y="temp", figsize=(12,6))


# In[19]:


flattened.write     .format("jdbc")     .mode("overwrite")     .option("url", "jdbc:postgresql://localhost:5432/postgres")     .option("dbtable", "public.test_table54")     .option("user", os.environ['LOCAL_POSTGRES_USER'])     .option("password", os.environ['LOCAL_POSTGRES_PASS'])     .option("driver", "org.postgresql.Driver")     .save()


# In[ ]:




